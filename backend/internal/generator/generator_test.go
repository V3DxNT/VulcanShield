package generator

import (
	"testing"
	"time"

	"github.com/vulcanshield/backend/internal/generator/scenarios"
	"github.com/vulcanshield/backend/internal/models"
)

var testPool = &models.EntityPool{
	Users: []models.UserProfile{
		{UserID: "C1001", TypicalMinAmount: 15.00, TypicalMaxAmount: 250.00, TrustScore: 85},
		{UserID: "C1002", TypicalMinAmount: 5.00, TypicalMaxAmount: 100.00, TrustScore: 60},
		{UserID: "C1003", TypicalMinAmount: 50.00, TypicalMaxAmount: 1500.00, TrustScore: 30},
	},
	DeviceIDs:   []string{"D204", "D205", "D206"},
	IPAddresses: []string{"IP-17", "IP-18", "IP-19"},
	MerchantIDs: []string{"M301", "M302", "M303"},
}

// TestDeterminism verifies that two generators with the same seed produce
// identical transaction sequences (AGENTS.md §22 — Demo Reliability Rule).
func TestDeterminism(t *testing.T) {
	seed := int64(42)

	gen1 := NewBaseGenerator(seed, testPool, &scenarios.NormalScenario{})
	gen2 := NewBaseGenerator(seed, testPool, &scenarios.NormalScenario{})

	for i := 0; i < 20; i++ {
		tx1 := gen1.Next(i, -1)
		tx2 := gen2.Next(i, -1)

		if tx1.TransactionID != tx2.TransactionID {
			t.Fatalf("idx %d: transaction IDs differ: %s != %s", i, tx1.TransactionID, tx2.TransactionID)
		}
		if tx1.UserID != tx2.UserID {
			t.Fatalf("idx %d: user IDs differ: %s != %s", i, tx1.UserID, tx2.UserID)
		}
		if tx1.Amount != tx2.Amount {
			t.Fatalf("idx %d: amounts differ: %f != %f", i, tx1.Amount, tx2.Amount)
		}
		if tx1.DeviceID != tx2.DeviceID {
			t.Fatalf("idx %d: device IDs differ: %s != %s", i, tx1.DeviceID, tx2.DeviceID)
		}
	}
}

// TestVelocityPinsUser verifies that velocity attack uses the same user/device/IP.
func TestVelocityPinsUser(t *testing.T) {
	gen := NewBaseGenerator(42, testPool, &scenarios.VelocityAttackScenario{})

	firstTx := gen.Next(0, -1)
	for i := 1; i < 15; i++ {
		tx := gen.Next(i, -1)
		if tx.UserID != firstTx.UserID {
			t.Errorf("velocity idx %d: user changed from %s to %s", i, firstTx.UserID, tx.UserID)
		}
		if tx.DeviceID != firstTx.DeviceID {
			t.Errorf("velocity idx %d: device changed from %s to %s", i, firstTx.DeviceID, tx.DeviceID)
		}
		if tx.IPAddress != firstTx.IPAddress {
			t.Errorf("velocity idx %d: IP changed from %s to %s", i, firstTx.IPAddress, tx.IPAddress)
		}
	}
}

// TestAmountAnomaly verifies that the amount is significantly above the user's typical max.
func TestAmountAnomaly(t *testing.T) {
	gen := NewBaseGenerator(42, testPool, &scenarios.AmountAnomalyScenario{})

	for i := 0; i < 10; i++ {
		tx := gen.Next(i, 0) // target user 0 = C1001, typical_max=250
		if tx.Amount < 250.0*5 {
			t.Errorf("amount anomaly idx %d: amount %f too low (expected >1250)", i, tx.Amount)
		}
	}
}

func TestGeneratedTimestampsProgressInRealTime(t *testing.T) {
	gen := NewBaseGenerator(42, testPool, &scenarios.NormalScenario{})
	prev := gen.Next(0, -1).Timestamp

	for i := 1; i < 8; i++ {
		tx := gen.Next(i, -1)
		if tx.Timestamp.Sub(prev) < 250*time.Millisecond {
			t.Fatalf("timestamp idx %d did not advance meaningfully: %s -> %s", i, prev.Format(time.RFC3339Nano), tx.Timestamp.Format(time.RFC3339Nano))
		}
		prev = tx.Timestamp
	}
}

// TestDeviceFarmSharesDevice verifies all transactions use the same device.
func TestDeviceFarmSharesDevice(t *testing.T) {
	gen := NewBaseGenerator(42, testPool, &scenarios.DeviceFarmScenario{})

	sharedDevice := testPool.DeviceIDs[len(testPool.DeviceIDs)-1] // D206
	for i := 0; i < 10; i++ {
		tx := gen.Next(i, -1)
		if tx.DeviceID != sharedDevice {
			t.Errorf("device farm idx %d: expected device %s, got %s", i, sharedDevice, tx.DeviceID)
		}
	}
	// Verify multiple users are cycled
	users := make(map[string]bool)
	for i := 0; i < 6; i++ {
		tx := gen.Next(i, -1)
		users[tx.UserID] = true
	}
	if len(users) < 2 {
		t.Errorf("device farm: expected multiple users, got %d", len(users))
	}
}

// TestIPAbuseSharesIP verifies all transactions use the same high-risk IP.
func TestIPAbuseSharesIP(t *testing.T) {
	gen := NewBaseGenerator(42, testPool, &scenarios.IPAbuseScenario{})

	sharedIP := testPool.IPAddresses[len(testPool.IPAddresses)-1] // IP-19
	for i := 0; i < 10; i++ {
		tx := gen.Next(i, -1)
		if tx.IPAddress != sharedIP {
			t.Errorf("ip abuse idx %d: expected IP %s, got %s", i, sharedIP, tx.IPAddress)
		}
	}
}

// TestAccountTakeoverUsesHighRisk verifies ATO uses high-risk device/IP and abnormal amount.
func TestAccountTakeoverUsesHighRisk(t *testing.T) {
	gen := NewBaseGenerator(42, testPool, &scenarios.AccountTakeoverScenario{})

	highRiskDevice := testPool.DeviceIDs[len(testPool.DeviceIDs)-1]
	highRiskIP := testPool.IPAddresses[len(testPool.IPAddresses)-1]

	for i := 0; i < 5; i++ {
		tx := gen.Next(i, 0) // target user 0 = C1001
		if tx.DeviceID != highRiskDevice {
			t.Errorf("ATO idx %d: expected high-risk device %s, got %s", i, highRiskDevice, tx.DeviceID)
		}
		if tx.IPAddress != highRiskIP {
			t.Errorf("ATO idx %d: expected high-risk IP %s, got %s", i, highRiskIP, tx.IPAddress)
		}
		// Amount should be well above typical max of 250
		if tx.Amount < 250.0*3 {
			t.Errorf("ATO idx %d: amount %f too low for takeover", i, tx.Amount)
		}
	}
}

// TestScenarioFor verifies the factory returns correct implementations.
func TestScenarioFor(t *testing.T) {
	tests := []struct {
		input    models.ScenarioType
		expected models.ScenarioType
	}{
		{models.ScenarioNormal, models.ScenarioNormal},
		{models.ScenarioVelocityAttack, models.ScenarioVelocityAttack},
		{models.ScenarioAccountTakeover, models.ScenarioAccountTakeover},
		{models.ScenarioDeviceFarm, models.ScenarioDeviceFarm},
		{models.ScenarioIPAbuse, models.ScenarioIPAbuse},
		{models.ScenarioAmountAnomaly, models.ScenarioAmountAnomaly},
		{"unknown_type", models.ScenarioNormal}, // unknown falls back to normal
	}

	for _, tc := range tests {
		s := ScenarioFor(tc.input)
		if s.Type() != tc.expected {
			t.Errorf("ScenarioFor(%q) = %q, want %q", tc.input, s.Type(), tc.expected)
		}
	}
}
