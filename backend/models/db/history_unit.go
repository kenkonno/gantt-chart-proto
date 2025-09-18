package db

import "time"

type HistoryUnit struct {
	Id         *int32 `gorm:"primaryKey"`
	HistoryId  int32  `gorm:"primaryKey"`
	Name       string
	FacilityId int32
	Order      int
	CreatedAt  time.Time
	UpdatedAt  int32
}

func (m *HistoryUnit) TableName() string {
	return "history_units"
}
