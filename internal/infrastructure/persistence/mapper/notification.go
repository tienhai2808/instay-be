package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
)

func ToNotificationEntity(nt *model.Notification) *entity.NotificationEntity {
	if nt == nil {
		return nil
	}

	return &entity.NotificationEntity{
		ID:            nt.ID,
		DepartmentID:  nt.DepartmentID,
		BookingRoomID: nt.BookingRoomID,
		ReceiverType:  string(nt.ReceiverType),
		ContentType:   string(nt.ContentType),
		ContentKey:    nt.ContentKey,
		ContentVI:     nt.ContentVI,
		ContentEN:     nt.ContentEN,
		IsRead:        nt.IsRead,
		CreatedAt:     nt.CreatedAt,
	}
}

func ToNotification(nt *entity.NotificationEntity) *model.Notification {
	if nt == nil {
		return nil
	}

	return &model.Notification{
		ID:            nt.ID,
		DepartmentID:  nt.DepartmentID,
		BookingRoomID: nt.BookingRoomID,
		ReceiverType:  model.NotificationReceiverType(nt.ReceiverType),
		ContentType:   model.NotificationContentType(nt.ContentType),
		ContentKey:    nt.ContentKey,
		ContentVI:     nt.ContentVI,
		ContentEN:     nt.ContentEN,
		IsRead:        nt.IsRead,
		CreatedAt:     nt.CreatedAt,
		Department:    ToDepartment(nt.Department),
		BookingRoom:   ToBookingRoom(nt.BookingRoom),
		Users:         ToUsers(nt.Users),
	}
}

func ToNotifications(nts []*entity.NotificationEntity) []*model.Notification {
	if len(nts) == 0 {
		return make([]*model.Notification, 0)
	}

	ntsRes := make([]*model.Notification, 0, len(nts))
	for _, nt := range nts {
		ntsRes = append(ntsRes, ToNotification(nt))
	}

	return ntsRes
}

func ToUserNotificationEntity(un *model.UserNotification) *entity.UserNotificationEntity {
	if un == nil {
		return nil
	}

	return &entity.UserNotificationEntity{
		UserID:         un.UserID,
		NotificationID: un.NotificationID,
	}
}

func ToUserNotificationEntities(uns []*model.UserNotification) []*entity.UserNotificationEntity {
	if len(uns) == 0 {
		return make([]*entity.UserNotificationEntity, 0)
	}

	unsRes := make([]*entity.UserNotificationEntity, 0, len(uns))
	for _, un := range uns {
		unsRes = append(unsRes, ToUserNotificationEntity(un))
	}

	return unsRes
}
