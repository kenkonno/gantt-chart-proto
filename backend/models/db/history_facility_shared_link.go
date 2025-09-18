package db

import (
	"time"
)

type HistoryFacilitySharedLink struct {
	Id         *int32 `gorm:"primaryKey"`
	HistoryId  int32  `gorm:"primaryKey"`
	FacilityId int32
	Uuid       string
	CreatedAt  time.Time
	UpdatedAt  int32
}

func (m *HistoryFacilitySharedLink) TableName() string {
	return "history_facility_shared_links"
}
