package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vulcanshield/backend/internal/models"
)

type GraphRepository interface {
	ListRelationships(ctx context.Context, limit int) ([]models.FraudRelationship, error)
	GetNeighbors(ctx context.Context, entityID string) ([]models.FraudRelationship, error)
	CreateRelationship(ctx context.Context, rel *models.FraudRelationship) error
}

type pgxGraphRepository struct {
	pool *pgxpool.Pool
}

func NewGraphRepository(pool *pgxpool.Pool) GraphRepository {
	return &pgxGraphRepository{pool: pool}
}

func (r *pgxGraphRepository) ListRelationships(ctx context.Context, limit int) ([]models.FraudRelationship, error) {
	const q = `
		SELECT relationship_id, source_type, source_id, target_type, target_id,
		       relationship_type, weight, fraud_linked, created_at
		FROM fraud_relationships
		ORDER BY created_at DESC
		LIMIT $1`

	rows, err := r.pool.Query(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rels []models.FraudRelationship
	for rows.Next() {
		var f models.FraudRelationship
		if err := rows.Scan(
			&f.RelationshipID, &f.SourceType, &f.SourceID, &f.TargetType, &f.TargetID,
			&f.RelationshipType, &f.Weight, &f.FraudLinked, &f.CreatedAt,
		); err != nil {
			return nil, err
		}
		rels = append(rels, f)
	}
	return rels, rows.Err()
}

func (r *pgxGraphRepository) GetNeighbors(ctx context.Context, entityID string) ([]models.FraudRelationship, error) {
	const q = `
		SELECT relationship_id, source_type, source_id, target_type, target_id,
		       relationship_type, weight, fraud_linked, created_at
		FROM fraud_relationships
		WHERE source_id = $1 OR target_id = $1
		ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, q, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rels []models.FraudRelationship
	for rows.Next() {
		var f models.FraudRelationship
		if err := rows.Scan(
			&f.RelationshipID, &f.SourceType, &f.SourceID, &f.TargetType, &f.TargetID,
			&f.RelationshipType, &f.Weight, &f.FraudLinked, &f.CreatedAt,
		); err != nil {
			return nil, err
		}
		rels = append(rels, f)
	}
	return rels, rows.Err()
}

func (r *pgxGraphRepository) CreateRelationship(ctx context.Context, rel *models.FraudRelationship) error {
	const q = `
		INSERT INTO fraud_relationships
		  (relationship_id, source_type, source_id, target_type, target_id, relationship_type, weight, fraud_linked, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (relationship_id) DO UPDATE SET weight = EXCLUDED.weight, fraud_linked = EXCLUDED.fraud_linked`

	_, err := r.pool.Exec(ctx, q,
		rel.RelationshipID, rel.SourceType, rel.SourceID, rel.TargetType, rel.TargetID,
		rel.RelationshipType, rel.Weight, rel.FraudLinked, rel.CreatedAt,
	)
	return err
}
