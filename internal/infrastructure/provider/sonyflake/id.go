package sonyflake

import (
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/sony/sonyflake/v2"
)

type IDProvider struct {
	gen *sonyflake.Sonyflake
}

func NewIDProvider(gen *sonyflake.Sonyflake) port.IDProvider {
	return &IDProvider{gen}
}

func (p *IDProvider) NextID() (int64, error) {
	return p.gen.NextID()
}
