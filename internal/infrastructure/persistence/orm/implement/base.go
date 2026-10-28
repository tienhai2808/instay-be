package implement

import (
	"context"
	"fmt"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"gorm.io/gorm"
)

type Preload struct {
	Relation string
	Scope    func(*gorm.DB) *gorm.DB
}

type BaseRepository struct {
	db *gorm.DB
}

func NewBaseRepository(db *gorm.DB) BaseRepository {
	return BaseRepository{db}
}

func (r *BaseRepository) getDB(ctx context.Context) *gorm.DB {
	if txDB, ok := ctx.Value(port.CtxTxKey{}).(*gorm.DB); ok {
		return txDB
	}

	return r.db.WithContext(ctx)
}

func (r *BaseRepository) schemaTable(schema, table string) string {
	return fmt.Sprintf(`"%s".%s`, schema, table)
}

func (r *BaseRepository) preloadWithSchema(schema, table string, selectFields ...string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Table(r.schemaTable(schema, table))
		if len(selectFields) > 0 {
			db = db.Select(selectFields)
		}
		return db
	}
}
