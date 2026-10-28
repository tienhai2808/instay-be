package implement

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
)

type schemaRepositoryImpl struct {
	BaseRepository
}

func NewSchemaRepository(baseRepo BaseRepository) repository.SchemeRepository {
	return &schemaRepositoryImpl{baseRepo}
}

func (r *schemaRepositoryImpl) FindAll(ctx context.Context) ([]string, error) {
	var schemas []string
	query := `
		SELECT nspname
		FROM pg_namespace
		WHERE nspname NOT LIKE 'pg_%'
		AND nspname NOT IN ('information_schema', 'public', 'atlas_schema_revisions')
		ORDER BY nspname
	`
	if err := r.db.WithContext(ctx).Raw(query).Scan(&schemas).Error; err != nil {
		return nil, err
	}

	return schemas, nil
}
