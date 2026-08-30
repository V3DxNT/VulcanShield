package scenarios

import (
	"math/rand"

	"github.com/vulcanshield/backend/internal/models"
)

// Scenario applies modifiers to a base transaction template produced by BaseGenerator.
// Each concrete scenario implements this interface.
type Scenario interface {
	// Type returns the canonical scenario name.
	Type() models.ScenarioType

	// Next produces the next transaction for the given sequence index.
	// rng is the seeded generator; pool contains all available synthetic entities.
	// targetUserIndex optionally pins a specific user (for attack scenarios); -1 means auto-select.
	Next(idx int, rng *rand.Rand, pool *models.EntityPool, targetUserIndex int) models.Transaction
}
