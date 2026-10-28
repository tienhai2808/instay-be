package model

type Currency string

const (
	CurrencyUSD Currency = "usd"
	CurrencyVND Currency = "vnd"
)

type Features struct{}

type Plan struct {
	ID           int64
	Name         string
	Price        float64
	Currency     Currency
	Features     Features
	Duration     string
	DurationDays int64
}
