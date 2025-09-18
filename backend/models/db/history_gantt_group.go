package db

import "time"

type HistoryGanttGroup struct {
	Id         *int32 `gorm:"primaryKey"`
	HistoryId  int32  `gorm:"primaryKey"`
	FacilityId int32
	UnitId     int32
	CreatedAt  time.Time
	UpdatedAt  int32
}

func (m *HistoryGanttGroup) TableName() string {
	return "history_gantt_groups"
}
