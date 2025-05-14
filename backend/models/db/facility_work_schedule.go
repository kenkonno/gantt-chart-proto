package db

import "time"

type FacilityWorkSchedule struct {
	Id         *int32 `gorm:"primaryKey;autoIncrement"`
	FacilityId int32  `gorm:"index:idx_facility_id"`
	Date       time.Time
	Type       string
	CreatedAt  time.Time
	UpdatedAt  int32
}
