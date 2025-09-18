package db

import "time"

type HistoryFacilityWorkSchedule struct {
	Id         *int32 `gorm:"primaryKey"`
	HistoryId  int32  `gorm:"primaryKey"`
	FacilityId int32  `gorm:"index:idx_facility_id"`
	Date       time.Time
	Type       string
	CreatedAt  time.Time
	UpdatedAt  int32
}

func (m *HistoryFacilityWorkSchedule) TableName() string {
	return "history_facility_work_schedules"
}
