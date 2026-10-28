package entity

const ServiceImageTableName = "service_images"

type ServiceImageEntity struct {
	ServiceID   int64  `gorm:"type:bigint;primaryKey" json:"service_id"`
	Key         string `gorm:"type:varchar(150);not null;uniqueIndex:service_images_key_key" json:"key"`
	IsThumbnail bool   `gorm:"type:boolean;not null" json:"is_thumbnail"`
	SortOrder   uint32 `gorm:"type:integer;primaryKey" json:"sort_order"`

	Service *ServiceEntity `gorm:"foreignKey:ServiceID;references:ID;constraint:fk_service_images_service,OnUpdate:CASCADE,OnDelete:CASCADE" json:"service"`
}

func (ServiceImageEntity) TableName() string {
	return ServiceImageTableName
}
