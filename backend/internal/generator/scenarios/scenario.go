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

func userDeviceID(pool *models.EntityPool, user models.UserProfile) string {
	if pool == nil {
		return ""
	}
	if pool.UserDeviceMap != nil {
		if deviceID, ok := pool.UserDeviceMap[user.UserID]; ok && deviceID != "" {
			return deviceID
		}
	}
	if len(pool.DeviceIDs) == 0 {
		return ""
	}
	if user.UserID == "" {
		return pool.DeviceIDs[0]
	}
	last := user.UserID[len(user.UserID)-1]
	idx := int(last)
	if last >= '0' && last <= '9' {
		idx = int(last - '0')
	}
	return pool.DeviceIDs[idx%len(pool.DeviceIDs)]
}

func userIP(pool *models.EntityPool, user models.UserProfile) string {
	if pool == nil {
		return ""
	}
	if pool.UserIPMap != nil {
		if ip, ok := pool.UserIPMap[user.UserID]; ok && ip != "" {
			return ip
		}
	}
	if len(pool.IPAddresses) == 0 {
		return ""
	}
	if user.UserID == "" {
		return pool.IPAddresses[0]
	}
	last := user.UserID[len(user.UserID)-1]
	idx := int(last)
	if last >= '0' && last <= '9' {
		idx = int(last - '0')
	}
	return pool.IPAddresses[idx%len(pool.IPAddresses)]
}
