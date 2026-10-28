package seeder

import (
	"fmt"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
	"gorm.io/gorm"
)

var autoMigrateModels = []any{
	&entity.DepartmentEntity{},
	&entity.PermissionEntity{},
	&entity.UserEntity{},
	&entity.RefreshTokenEntity{},
	&entity.DeviceTokenEntity{},
	&entity.ServiceTypeEntity{},
	&entity.ServiceEntity{},
	&entity.ServiceImageEntity{},
	&entity.SupportCategoryEntity{},
	&entity.SourceEntity{},
	&entity.BranchEntity{},
	&entity.UserBranchEntity{},
	&entity.BranchServiceEntity{},
	&entity.RoomTypeEntity{},
	&entity.BranchRoomTypeEntity{},
	&entity.FloorEntity{},
	&entity.RoomEntity{},
	&entity.BookingEntity{},
	&entity.BookingRoomEntity{},
	&entity.SupportTicketEntity{},
	&entity.OrderServiceEntity{},
	&entity.NotificationEntity{},
	&entity.GuestEntity{},
	&entity.CustomerEntity{},
	&entity.ReviewEntity{},
	&entity.ChatEntity{},
	&entity.PricingTagEntity{},
	&entity.PricingCalendarEntity{},
	&entity.RoomTypePriceEntity{},
	&entity.DailyRoomTypeRateEntity{},
	&entity.CustomerEntity{},
	&entity.EventEntity{},
	&entity.BranchEventEntity{},
}

type Seeder struct {
	logPro         port.LogProvider
	db             *gorm.DB
	idPro          port.IDProvider
	userRepo       repository.UserRepository
	schemaRepo     repository.SchemeRepository
	permissionRepo repository.PermissionRepository
}

func NewSeeder(
	logPro port.LogProvider,
	db *gorm.DB,
	idPro port.IDProvider,
	userRepo repository.UserRepository,
	schemaRepo repository.SchemeRepository,
	permissionRepo repository.PermissionRepository,
) *Seeder {
	return &Seeder{
		logPro,
		db,
		idPro,
		userRepo,
		schemaRepo,
		permissionRepo,
	}
}

func (s *Seeder) Start(seedName string) error {
	switch seedName {
	case "admin":
		return s.adminSeeder()
	case "tenant":
		return s.tenantSeeder()
	case "migrate":
		return s.migrateSeeder()
	default:
		return fmt.Errorf("unknown seed: %s", seedName)
	}
}
