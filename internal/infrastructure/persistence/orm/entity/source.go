package entity

import "github.com/InstaySystem/is_v2-be/internal/domain/model"

const SourceTableName = "sources"

type SourceEntity struct {
	BaseEntity
	AuditEntity

	Name string `gorm:"type:varchar(50);not null;uniqueIndex:sources_name_key" json:"name"`

	Bookings []*BookingEntity `gorm:"foreignKey:SourceID;references:ID;constraint:fk_bookings_source,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"bookings"`
}

func (SourceEntity) TableName() string {
	return SourceTableName
}

func DefaultSources() []SourceEntity {
	return []SourceEntity{
		{
			Name: model.SourceOffline,
		},
	}
}
