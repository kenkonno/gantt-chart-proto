package db

import "time"

type HistoryHoliday struct {
	Id        *int32 `gorm:"primaryKey"`
	HistoryId int32  `gorm:"primaryKey"`
	Name      string
	Date      time.Time
	CreatedAt time.Time
	UpdatedAt int32
}

func (m *HistoryHoliday) TableName() string {
	return "history_holidays"
}
