package implement

import (
	"context"
	"fmt"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/mapper"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
	"gorm.io/gorm"
)

type notificationRepositoryImpl struct {
	BaseRepository
}

func NewNotificationRepository(baseRepo BaseRepository) repository.NotificationRepository {
	return &notificationRepositoryImpl{baseRepo}
}

func (r *notificationRepositoryImpl) Create(ctx context.Context, schema string, notification *model.Notification) error {
	notificationEntity := mapper.ToNotificationEntity(notification)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.NotificationTableName)).
		Create(notificationEntity).Error
}

func (r *notificationRepositoryImpl) CountUnreadByUserIDAndDepartmentID(ctx context.Context, schema string, userID, departmentID int64, query model.NotificationQuery) (int64, error) {
	var count int64

	notiTable := r.schemaTable(schema, entity.NotificationTableName)
	userNotiTable := r.schemaTable(schema, entity.UserNotificationTableName)

	db := r.getDB(ctx).
		Table(notiTable+" AS n").
		Joins(fmt.Sprintf(`LEFT JOIN %s un ON un.notification_id = n.id AND un.user_id = ?`, userNotiTable), userID).
		Where("n.receiver_type = ?", model.NotificationReceiverTypeStaff).
		Where("n.department_id = ?", departmentID).
		Where("un.notification_id IS NULL")

	db = r.applyStaffBranchFilter(db, schema, query)

	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *notificationRepositoryImpl) CountUnreadByBookingRoomID(ctx context.Context, schema string, bookingRoomID int64) (int64, error) {
	var count int64

	notiTable := r.schemaTable(schema, entity.NotificationTableName)

	if err := r.getDB(ctx).
		Table(notiTable+" AS n").
		Where("n.receiver_type = ?", model.NotificationReceiverTypeGuest).
		Where("n.booking_room_id = ?", bookingRoomID).
		Where("n.is_read = ?", false).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *notificationRepositoryImpl) FindAllByUserIDAndDepartmentIDPaginated(ctx context.Context, schema string, userID, departmentID int64, query model.NotificationQuery) ([]*model.Notification, int64, error) {
	var notificationEntities []*entity.NotificationEntity
	var total int64

	notiTable := r.schemaTable(schema, entity.NotificationTableName)
	userNotiTable := r.schemaTable(schema, entity.UserNotificationTableName)

	db := r.getDB(ctx).
		Table(notiTable+" AS n").
		Model(&entity.NotificationEntity{}).
		Where("n.receiver_type = ?", model.NotificationReceiverTypeStaff).
		Where("n.department_id = ?", departmentID)

	db = r.applyStaffBranchFilter(db, schema, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.Notification{}, 0, nil
	}

	db = db.Session(&gorm.Session{})

	offset := (query.Page - 1) * query.Limit

	if err := db.
		Select(
			fmt.Sprintf(
				`n.id, n.content_type, n.content_key, n.content_vi, n.content_en, n.created_at,
				EXISTS(SELECT 1 FROM %s un WHERE un.notification_id = n.id AND un.user_id = ?) AS is_read`,
				userNotiTable,
			),
			userID,
		).
		Order("n.created_at DESC").
		Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&notificationEntities).Error; err != nil {
		return nil, 0, err
	}

	return mapper.ToNotifications(notificationEntities), total, nil
}

func (r *notificationRepositoryImpl) FindAllByBookingRoomIDPaginated(ctx context.Context, schema string, bookingRoomID int64, query model.NotificationQuery) ([]*model.Notification, int64, error) {
	var notificationEntities []*entity.NotificationEntity
	var total int64

	notiTable := r.schemaTable(schema, entity.NotificationTableName)

	db := r.getDB(ctx).
		Table(notiTable+" AS n").
		Model(&entity.NotificationEntity{}).
		Where("n.receiver_type = ?", model.NotificationReceiverTypeGuest).
		Where("n.booking_room_id = ?", bookingRoomID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.Notification{}, 0, nil
	}

	db = db.Session(&gorm.Session{})

	offset := (query.Page - 1) * query.Limit

	if err := db.
		Select("n.id, n.content_type, n.content_key, n.content_vi, n.content_en, n.is_read, n.created_at").
		Order("n.created_at DESC").
		Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&notificationEntities).Error; err != nil {
		return nil, 0, err
	}

	return mapper.ToNotifications(notificationEntities), total, nil
}

func (r *notificationRepositoryImpl) UpdateReadByBookingRoomIDAndIDs(ctx context.Context, schema string, bookingRoomID int64, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	notiTable := r.schemaTable(schema, entity.NotificationTableName)

	return r.getDB(ctx).
		Table(notiTable).
		Where("receiver_type = ?", model.NotificationReceiverTypeGuest).
		Where("booking_room_id = ?", bookingRoomID).
		Where("id IN ?", ids).
		Update("is_read", true).Error
}

func (r *notificationRepositoryImpl) applyStaffBranchFilter(db *gorm.DB, schema string, query model.NotificationQuery) *gorm.DB {
	if len(query.BranchIDs) == 0 {
		return db
	}

	bookingRoomTable := r.schemaTable(schema, entity.BookingRoomTableName)
	roomTable := r.schemaTable(schema, entity.RoomTableName)
	floorTable := r.schemaTable(schema, entity.FloorTableName)

	return db.
		Joins(fmt.Sprintf(`INNER JOIN %s br_scope ON br_scope.id = n.booking_room_id`, bookingRoomTable)).
		Joins(fmt.Sprintf(`INNER JOIN %s r_scope ON r_scope.id = br_scope.room_id`, roomTable)).
		Joins(fmt.Sprintf(`INNER JOIN %s f_scope ON f_scope.id = r_scope.floor_id`, floorTable)).
		Where("f_scope.branch_id IN ?", query.BranchIDs)
}
