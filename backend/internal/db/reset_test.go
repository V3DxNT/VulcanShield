package db

import (
	"strings"
	"testing"
)

func TestDemoResetSQL_LeavesSeedDataIntact(t *testing.T) {
	sql := DemoResetSQL()

	for _, keep := range runtimeResetTables {
		if !strings.Contains(sql, keep) {
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
		"fraud_relationships",
		"rag_documents",
		"embedding_records",
	} {
		if strings.Contains(sql, drop) {
			t.Fatalf("demo reset SQL should not clear seed table %q: %s", drop, sql)
		}
	}
}
