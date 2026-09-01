package generator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/vulcanshield/backend/internal/challenge"
	"github.com/vulcanshield/backend/internal/db/repository"
	"github.com/vulcanshield/backend/internal/features"
	"github.com/vulcanshield/backend/internal/graph"
	"github.com/vulcanshield/backend/internal/kafka"
	"github.com/vulcanshield/backend/internal/mlclient"
	"github.com/vulcanshield/backend/internal/models"
	"github.com/vulcanshield/backend/internal/policy"
	appredis "github.com/vulcanshield/backend/internal/redis"
	"github.com/vulcanshield/backend/internal/risk"
	appws "github.com/vulcanshield/backend/internal/websocket"
)


var ErrScenarioRunning = errors.New("a scenario is already running")


type Engine struct {
	mu             sync.Mutex
	active         *models.ScenarioRun
	cancelFn       context.CancelFunc
	done           chan struct{}
	log            *slog.Logger
	txRepo         repository.TransactionRepository
	scRepo         repository.ScenarioRepository
	entityRepo     repository.EntityRepository
	riskRepo       repository.RiskRepository
	policyRepo     repository.PolicyRepository
	challengeRepo  repository.ChallengeRepository
	userRepo       repository.UserRepository
	kafka          *kafka.Producer
	velocityEngine *appredis.VelocityEngine
	featureBuilder *features.FeatureBuilder
	mlClient       *mlclient.Client
	riskEvaluator  *risk.Evaluator
	policyEngine   *policy.Engine
	graphEngine    *graph.Engine
	otpService     *challenge.Service
	wsHub          *appws.Hub
}


func NewEngine(
	log *slog.Logger,
	txRepo repository.TransactionRepository,
	scRepo repository.ScenarioRepository,
	entityRepo repository.EntityRepository,
	riskRepo repository.RiskRepository,
	policyRepo repository.PolicyRepository,
	challengeRepo repository.ChallengeRepository,
	userRepo repository.UserRepository,
	kafkaProducer *kafka.Producer,
	velocityEngine *appredis.VelocityEngine,
	mlClient *mlclient.Client,
	otpService *challenge.Service,
	graphEngine *graph.Engine,
	wsHub *appws.Hub,
) *Engine {
	return &Engine{
		log:            log,
		txRepo:         txRepo,
		scRepo:         scRepo,
		entityRepo:     entityRepo,
		riskRepo:       riskRepo,
		policyRepo:     policyRepo,
		challengeRepo:  challengeRepo,
		userRepo:       userRepo,
		kafka:          kafkaProducer,
		velocityEngine: velocityEngine,
		featureBuilder: features.NewFeatureBuilder(),
		mlClient:       mlClient,
		riskEvaluator:  risk.NewEvaluator(),
		policyEngine:   policy.NewEngine(),
		graphEngine:    graphEngine,
		otpService:     otpService,
		wsHub:          wsHub,
	}
}



func (e *Engine) Start(ctx context.Context, req models.ScenarioStartRequest) (*models.ScenarioRun, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.active != nil && e.active.Status == models.ScenarioRunning {
		return nil, ErrScenarioRunning
	}

	
	if req.Transactions <= 0 {
		req.Transactions = 100
	}
	if req.IntervalMS <= 0 {
		req.IntervalMS = 1000
	}
	if req.Seed == 0 {
		req.Seed = 42
	}

	
	if !validScenario(req.Scenario) {
		req.Scenario = models.ScenarioNormal
	}

	run := &models.ScenarioRun{
		ScenarioID:       fmt.Sprintf("SC-%d", time.Now().UnixMilli()),
		ScenarioType:     req.Scenario,
		TransactionCount: req.Transactions,
		GeneratedCount:   0,
		Seed:             req.Seed,
		Status:           models.ScenarioRunning,
		StartedAt:        time.Now().UTC(),
	}

	if err := e.scRepo.Create(ctx, run); err != nil {
		return nil, fmt.Errorf("persisting scenario run: %w", err)
	}

	runCtx, cancel := context.WithCancel(context.Background())
	e.active = run
	e.cancelFn = cancel
	e.done = make(chan struct{})

	go e.runLoop(runCtx, run, req)
	return run, nil
}



func (e *Engine) Stop(ctx context.Context) (*models.ScenarioRun, error) {
	e.mu.Lock()
	if e.active == nil || e.active.Status != models.ScenarioRunning {
		e.mu.Unlock()
		return e.active, nil
	}
	cancel := e.cancelFn
	done := e.done
	e.mu.Unlock()

	cancel()
	<-done 

	e.mu.Lock()
	defer e.mu.Unlock()
	return e.active, nil
}


func (e *Engine) GetStatus() *models.ScenarioRun {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.active
}


func (e *Engine) runLoop(ctx context.Context, run *models.ScenarioRun, req models.ScenarioStartRequest) {
	defer close(e.done)

	
	pool, err := e.entityRepo.LoadPool(ctx)
	if err != nil {
		e.log.Error("generator: failed to load entity pool", "error", err, "scenario_id", run.ScenarioID)
		e.finalize(run, models.ScenarioStopped)
		return
	}
	if len(pool.Users) == 0 || len(pool.DeviceIDs) == 0 {
		e.log.Error("generator: entity pool is empty — run seed data", "scenario_id", run.ScenarioID)
		e.finalize(run, models.ScenarioStopped)
		return
	}

	userMap := make(map[string]*models.UserProfile)
	for i := range pool.Users {
		userMap[pool.Users[i].UserID] = &pool.Users[i]
	}

	scene := ScenarioFor(req.Scenario)
	gen := NewBaseGenerator(req.Seed, pool, scene)

	targetUserIndex := -1
	if req.CustomerID != "" {
		for i, u := range pool.Users {
			if u.UserID == req.CustomerID {
				targetUserIndex = i
				break
			}
		}
	}

	ticker := time.NewTicker(time.Duration(req.IntervalMS) * time.Millisecond)
	defer ticker.Stop()

	generated := 0
	for generated < req.Transactions {
		select {
		case <-ctx.Done():
			e.log.Info("generator: stopped by request", "scenario_id", run.ScenarioID, "generated", generated)
			e.finalizeWithCount(run, models.ScenarioStopped, generated)
			return
		case <-ticker.C:
			tx := gen.Next(generated, targetUserIndex)

			
			if e.velocityEngine != nil {
				_ = e.velocityEngine.RecordTransaction(ctx, &tx)
			}
			velocitySignals, _ := e.velocityEngine.GetVelocitySignals(ctx, &tx)

			
			isEmulator := tx.DeviceID == pool.DeviceIDs[len(pool.DeviceIDs)-1]
			isVPN := tx.IPAddress == pool.IPAddresses[len(pool.IPAddresses)-1]

			
			userProf := userMap[tx.UserID]
			predictReq := e.featureBuilder.BuildVector(&tx, userProf, velocitySignals, isEmulator, isVPN)

			
			var mlResp *mlclient.PredictResponse
			if e.mlClient != nil {
				var err error
				mlResp, err = e.mlClient.Predict(ctx, predictReq)
				if err != nil {
					e.log.Warn("generator: ML predict failed", "error", err, "transaction_id", tx.TransactionID)
				}
			}

			
			riskAssessment := e.riskEvaluator.Evaluate(&tx, mlResp, velocitySignals)

			
			var userThresholds *models.User
			if userProf != nil {
				userThresholds = &models.User{
					ChallengeThreshold: userProf.ChallengeThreshold,
					BlockThreshold:     userProf.BlockThreshold,
				}
			}
			policyDecision, updatedStatus := e.policyEngine.Evaluate(&tx, riskAssessment, userProf, userThresholds)
			tx.Status = updatedStatus

			
			if err := e.txRepo.Create(ctx, &tx); err != nil {
				if ctx.Err() != nil {
					e.finalizeWithCount(run, models.ScenarioStopped, generated)
					return
				}
				e.log.Error("generator: failed to persist transaction", "error", err,
					"transaction_id", tx.TransactionID, "scenario_id", run.ScenarioID)
				e.finalizeWithCount(run, models.ScenarioStopped, generated)
				return
			}

			
			if e.riskRepo != nil && riskAssessment != nil {
				if err := e.riskRepo.Create(ctx, riskAssessment); err != nil {
					e.log.Warn("generator: failed to persist risk assessment", "error", err, "transaction_id", tx.TransactionID)
				}
			}
			if e.policyRepo != nil && policyDecision != nil {
				if err := e.policyRepo.Create(ctx, policyDecision); err != nil {
					e.log.Warn("generator: failed to persist policy decision", "error", err, "transaction_id", tx.TransactionID)
				}
			}

			
			if policyDecision.Decision == models.DecisionChallenge && e.otpService != nil {
				otpChallenge, _, err := e.otpService.GenerateChallenge(ctx, tx.TransactionID)
				if err != nil {
					e.log.Warn("generator: OTP generate failed", "error", err)
				} else if e.challengeRepo != nil && otpChallenge != nil {
					if err := e.challengeRepo.Create(ctx, otpChallenge); err != nil {
						e.log.Warn("generator: failed to persist OTP challenge", "error", err, "transaction_id", tx.TransactionID)
					} else {
						e.settleDemoOTP(ctx, &tx, otpChallenge, generated, policyDecision)
					}
				}
			}

			if e.graphEngine != nil {
				if err := e.graphEngine.RecordTransactionEdges(ctx, &tx, isEmulator, isVPN); err != nil {
					e.log.Warn("generator: failed to record graph edges", "error", err, "transaction_id", tx.TransactionID)
				}
			}

			
			auditDetails := map[string]any{
				"scenario_id":   run.ScenarioID,
				"scenario_type": string(run.ScenarioType),
				"amount":        tx.Amount,
				"user_id":       tx.UserID,
				"risk_score":    riskAssessment.RiskScore,
			}
			_ = e.scRepo.InsertAuditEvent(ctx, tx.TransactionID, "TRANSACTION_CREATED", auditDetails)
			_ = e.scRepo.InsertAuditEvent(ctx, tx.TransactionID, "RISK_CALCULATED", map[string]any{
				"risk_score":        riskAssessment.RiskScore,
				"fraud_probability": riskAssessment.FraudProbability,
				"anomaly_score":     riskAssessment.AnomalyScore,
			})

			
			if e.kafka != nil {
				txPayload, _ := json.Marshal(tx)
				if err := e.kafka.Produce(ctx, kafka.TopicTransactionCreated, tx.TransactionID, txPayload); err != nil {
					e.log.Warn("generator: kafka produce failed", "topic", kafka.TopicTransactionCreated, "error", err)
				}

				riskPayload, _ := json.Marshal(riskAssessment)
				if err := e.kafka.Produce(ctx, kafka.TopicRiskEvaluated, tx.TransactionID, riskPayload); err != nil {
					e.log.Warn("generator: kafka produce failed", "topic", kafka.TopicRiskEvaluated, "error", err)
				}

				if policyDecision != nil {
					decisionPayload, _ := json.Marshal(policyDecision)
					if err := e.kafka.Produce(ctx, kafka.TopicTransactionDecisioned, tx.TransactionID, decisionPayload); err != nil {
						e.log.Warn("generator: kafka produce failed", "topic", kafka.TopicTransactionDecisioned, "error", err)
					}
				}
			}

			
			if e.wsHub != nil {
				e.wsHub.Broadcast("transaction_created", tx)
				e.wsHub.Broadcast("risk_updated", riskAssessment)
				e.wsHub.Broadcast("decision_created", policyDecision)
			}

			generated++
			e.mu.Lock()
			run.GeneratedCount = generated
			e.mu.Unlock()

			e.log.Debug("generator: transaction pipeline completed",
				"transaction_id", tx.TransactionID,
				"amount", tx.Amount,
				"risk_score", riskAssessment.RiskScore,
				"scenario_id", run.ScenarioID,
			)
		}
	}

	e.log.Info("generator: scenario completed", "scenario_id", run.ScenarioID, "generated", generated)
	e.finalizeWithCount(run, models.ScenarioCompleted, generated)
}



func (e *Engine) settleDemoOTP(
	ctx context.Context,
	tx *models.Transaction,
	otpChallenge *models.OTPChallenge,
	generated int,
	prior *models.PolicyDecision,
) {
	now := time.Now().UTC()
	approve := generated%2 == 0
	riskScore, challengeTh, blockTh := 0, 65, 85
	if prior != nil {
		riskScore = prior.RiskScore
		challengeTh = prior.ChallengeThreshold
		blockTh = prior.BlockThreshold
	}

	if approve {
		otpChallenge.Status = models.ChallengeVerified
		otpChallenge.Attempts = 1
		otpChallenge.VerifiedAt = &now
		_ = e.challengeRepo.Update(ctx, otpChallenge)
		tx.Status = models.StatusApproved
		_ = e.txRepo.UpdateStatus(ctx, tx.TransactionID, models.StatusApproved)
		finalRiskScore := e.recordVerifiedOTPAssessment(ctx, tx.TransactionID, riskScore, challengeTh)
		if e.policyRepo != nil {
			_ = e.policyRepo.Create(ctx, &models.PolicyDecision{
				DecisionID:         "PD-OTP-OK-" + tx.TransactionID,
				TransactionID:      tx.TransactionID,
				Decision:           models.DecisionAllow,
				RiskScore:          finalRiskScore,
				ChallengeThreshold: challengeTh,
				BlockThreshold:     blockTh,
				PolicyVersion:      "v1.0",
				Reason:             "Step-up OTP verified — customer confirmed; transaction approved",
				RulesTriggered:     []string{"RULE_OTP_VERIFIED"},
				CreatedAt:          now,
			})
		}
		e.log.Info("generator: demo OTP verified", "transaction_id", tx.TransactionID)
		return
	}

	otpChallenge.Status = models.ChallengeFailed
	otpChallenge.Attempts = otpChallenge.MaxAttempts
	_ = e.challengeRepo.Update(ctx, otpChallenge)
	tx.Status = models.StatusBlocked
	_ = e.txRepo.UpdateStatus(ctx, tx.TransactionID, models.StatusBlocked)
	if e.policyRepo != nil {
		_ = e.policyRepo.Create(ctx, &models.PolicyDecision{
			DecisionID:         "PD-OTP-FAIL-" + tx.TransactionID,
			TransactionID:      tx.TransactionID,
			Decision:           models.DecisionBlock,
			RiskScore:          riskScore,
			ChallengeThreshold: challengeTh,
			BlockThreshold:     blockTh,
			PolicyVersion:      "v1.0",
			Reason:             "Step-up OTP failed — max attempts exceeded; transaction blocked",
			RulesTriggered:     []string{"RULE_OTP_FAILED"},
			CreatedAt:          now,
		})
	}
	e.log.Info("generator: demo OTP rejected", "transaction_id", tx.TransactionID)
}




func (e *Engine) recordVerifiedOTPAssessment(ctx context.Context, transactionID string, priorScore, challengeThreshold int) int {
	finalScore := priorScore - 30
	maxAllowed := challengeThreshold - 1
	if finalScore > maxAllowed {
		finalScore = maxAllowed
	}
	if finalScore < 0 {
		finalScore = 0
	}
	if e.riskRepo == nil {
		return finalScore
	}
	prior, err := e.riskRepo.GetByTransactionID(ctx, transactionID)
	if err != nil || prior == nil {
		return finalScore
	}
	snapshot := make(map[string]any, len(prior.FeatureSnapshot)+3)
	for key, value := range prior.FeatureSnapshot {
		snapshot[key] = value
	}
	snapshot["challenge_result"] = "OTP_VERIFIED"
	snapshot["prior_risk_score"] = prior.RiskScore
	snapshot["risk_adjustment"] = finalScore - prior.RiskScore
	_ = e.riskRepo.Create(ctx, &models.RiskAssessment{
		AssessmentID:        "RA-OTP-OK-" + transactionID,
		TransactionID:       transactionID,
		FraudProbability:    prior.FraudProbability,
		AnomalyScore:        prior.AnomalyScore,
		FraudModelVersion:   prior.FraudModelVersion,
		AnomalyModelVersion: prior.AnomalyModelVersion,
		RiskScore:           finalScore,
		FeatureSnapshot:     snapshot,
		CreatedAt:           time.Now().UTC(),
	})
	return finalScore
}


func (e *Engine) finalize(run *models.ScenarioRun, status models.ScenarioStatus) {
	e.finalizeWithCount(run, status, run.GeneratedCount)
}

func (e *Engine) finalizeWithCount(run *models.ScenarioRun, status models.ScenarioStatus, count int) {
	now := time.Now().UTC()
	e.mu.Lock()
	run.Status = status
	run.GeneratedCount = count
	run.EndedAt = &now
	e.mu.Unlock()

	updateCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.scRepo.UpdateStatus(updateCtx, run.ScenarioID, status, &now); err != nil {
		e.log.Warn("generator: failed to update scenario status", "error", err, "scenario_id", run.ScenarioID)
	}
}


func validScenario(t models.ScenarioType) bool {
	switch t {
	case models.ScenarioNormal,
		models.ScenarioVelocityAttack,
		models.ScenarioAccountTakeover,
		models.ScenarioDeviceFarm,
		models.ScenarioIPAbuse,
		models.ScenarioAmountAnomaly:
		return true
	}
	return false
}
