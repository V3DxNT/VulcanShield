package db

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

var runtimeResetTables = []string{
	"audit_events",
	"investigation_evidence",
	"investigations",
	"otp_challenges",
	"policy_decisions",
	"risk_assessments",
	"transactions",
	"scenarios",
	"fraud_relationships",
}



func DemoResetSQL() string {
	return "TRUNCATE TABLE " + strings.Join(runtimeResetTables, ", ") + " RESTART IDENTITY CASCADE;"
}



func ResetRuntimeState(ctx context.Context, pool *pgxpool.Pool, flushRedis func(context.Context) error) error {
	if pool == nil {
		return nil
	}

	if _, err := pool.Exec(ctx, DemoResetSQL()); err != nil {
		return err
	}

	if flushRedis != nil {
		if err := flushRedis(ctx); err != nil {
			return err
		}
	}

	return nil
}
