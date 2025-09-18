package db

import (
	"time"
)

type Facility struct {
	Id              *int32 `gorm:"primaryKey;autoIncrement"`
	Name            string
	TermFrom        time.Time
	TermTo          time.Time
	Order           int
	Status          string
	Type            string
	FreeText *string
	ShipmentDueDate time.Time
	CreatedAt       time.Time
	UpdatedAt       int32
}
