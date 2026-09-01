package db

import (
	"strings"
	"testing"
)

func containsTableName(sql, name string) bool {
	return strings.Contains(sql, name+",") ||
		strings.Contains(sql, ", "+name) ||
		strings.Contains(sql, " "+name+" RESTART") ||
		strings.Contains(sql, "TRUNCATE TABLE "+name)
}

func TestDemoResetSQL_LeavesSeedDataIntact(t *testing.T) {
	sql := DemoResetSQL()

	for _, keep := range []string{
		"audit_events",
		"investigation_evidence",
		"investigations",
		"otp_challenges",
		"policy_decisions",
		"risk_assessments",
		"transactions",
		"scenarios",
		"fraud_relationships",
	} {
		if !containsTableName(sql, keep) {
			t.Fatalf("demo reset SQL missing table %q: %s", keep, sql)
		}
	}

	for _, drop := range []string{
		"users",
		"devices",
		"ips",
		"merchants",
		"user_devices",
		"user_ips",
		"rag_documents",
		"embedding_records",
	} {
		if containsTableName(sql, drop) {
			t.Fatalf("demo reset SQL should not clear seed table %q: %s", drop, sql)
		}
	}
}
