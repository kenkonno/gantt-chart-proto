package db

import (
	"time"
)

type Milestone struct {
	Id              *int32 `gorm:"primaryKey;autoIncrement"`
	FacilityId      int32
	Date            time.Time
	Description     string
	Order           int
	CreatedAt       time.Time
	UpdatedAt       int32
}
