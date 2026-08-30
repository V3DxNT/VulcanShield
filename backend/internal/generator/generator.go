package generator

import (
	"math/rand"

	"github.com/vulcanshield/backend/internal/generator/scenarios"
	"github.com/vulcanshield/backend/internal/models"
)

// BaseGenerator holds the seeded RNG and entity pool, and delegates
// transaction construction to a Scenario implementation.
// The generator is intentionally separated from persistence and event
// publication (AGENTS.md §4 constraint: generator logic separated from
// persistence and Kafka).
type BaseGenerator struct {
	rng   *rand.Rand
	pool  *models.EntityPool
	scene scenarios.Scenario
}

// NewBaseGenerator creates a deterministic generator seeded with seed.
// The same seed will always produce the same sequence of transactions
// for reproducible hackathon demonstrations (AGENTS.md §22).
func NewBaseGenerator(seed int64, pool *models.EntityPool, scene scenarios.Scenario) *BaseGenerator {
	return &BaseGenerator{
		rng:   rand.New(rand.NewSource(seed)), //nolint:gosec — seeded for reproducibility not security
		pool:  pool,
		scene: scene,
	}
}

// Next returns the transaction at position idx in the sequence.
// targetUserIndex < 0 means "auto-select per-scenario logic".
func (g *BaseGenerator) Next(idx int, targetUserIndex int) models.Transaction {
	return g.scene.Next(idx, g.rng, g.pool, targetUserIndex)
}

// ScenarioFor returns the Scenario implementation for a given type.
// Returns a NormalScenario for unrecognised types.
func ScenarioFor(t models.ScenarioType) scenarios.Scenario {
	switch t {
	case models.ScenarioVelocityAttack:
		return &scenarios.VelocityAttackScenario{}
	case models.ScenarioAccountTakeover:
		return &scenarios.AccountTakeoverScenario{}
	case models.ScenarioDeviceFarm:
		return &scenarios.DeviceFarmScenario{}
	case models.ScenarioIPAbuse:
		return &scenarios.IPAbuseScenario{}
	case models.ScenarioAmountAnomaly:
		return &scenarios.AmountAnomalyScenario{}
	default:
		return &scenarios.NormalScenario{}
	}
}
