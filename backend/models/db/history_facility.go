package db

import (
	"time"
)

type HistoryFacility struct {
	Id         *int32 `gorm:"primaryKey;autoIncrement"`
	FacilityId int32
	Name       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (m *HistoryFacility) TableName() string {
	return "facility_histories"
}
