package db

import (
	"time"
)

type HistoryMilestone struct {
	Id          *int32 `gorm:"primaryKey"`
	HistoryId   int32  `gorm:"primaryKey"`
	FacilityId  int32
	Date        time.Time
	Description string
	Order       int
	CreatedAt   time.Time
	UpdatedAt   int32
}

func (m *HistoryMilestone) TableName() string {
	return "history_milestones"
}
