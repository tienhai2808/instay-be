package implement

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/mapper"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
)

type serviceImageRepositoryImpl struct {
	BaseRepository
}

func NewServiceImageRepository(baseRepo BaseRepository) repository.ServiceImageRepository {
	return &serviceImageRepositoryImpl{baseRepo}
}

func (r *serviceImageRepositoryImpl) CreateBatch(ctx context.Context, schema string, serviceImages []*model.ServiceImage) error {
	svcImageEntities := mapper.ToServiceImageEntities(serviceImages)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ServiceImageTableName)).
		Create(svcImageEntities).Error
}

func (r *serviceImageRepositoryImpl) DeleteByServiceID(ctx context.Context, schema string, serviceID int64) error {
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ServiceImageTableName)).
		Where("service_id = ?", serviceID).
		Delete(&entity.ServiceImageEntity{}).Error
}
