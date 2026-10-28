package initialization

import (
	"time"

	"github.com/sony/sonyflake/v2"
)

type IDGen struct {
	sf *sonyflake.Sonyflake
}

func InitIDGen() (*IDGen, error) {
	st := sonyflake.Settings{
		StartTime: time.Date(2026, 1, 7, 0, 0, 0, 0, time.UTC),
		MachineID: func() (int, error) {
			return 2, nil
		},
	}

	sf, err := sonyflake.New(st)
	if err != nil {
		return nil, err
	}

	return &IDGen{
		sf,
	}, nil
}

func (i *IDGen) Generator() *sonyflake.Sonyflake {
	return i.sf
}
