package generator

import (
	"math/rand"

	"github.com/vulcanshield/backend/internal/generator/scenarios"
	"github.com/vulcanshield/backend/internal/models"
)






type BaseGenerator struct {
	rng   *rand.Rand
	pool  *models.EntityPool
	scene scenarios.Scenario
}




func NewBaseGenerator(seed int64, pool *models.EntityPool, scene scenarios.Scenario) *BaseGenerator {
	return &BaseGenerator{
		rng:   rand.New(rand.NewSource(seed)), 
		pool:  pool,
		scene: scene,
	}
}



func (g *BaseGenerator) Next(idx int, targetUserIndex int) models.Transaction {
	return g.scene.Next(idx, g.rng, g.pool, targetUserIndex)
}



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
