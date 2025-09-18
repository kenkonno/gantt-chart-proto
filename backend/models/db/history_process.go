package db

import "time"

type HistoryProcess struct {
	Id        *int32 `gorm:"primaryKey"`
	HistoryId int32  `gorm:"primaryKey"`
	Name      string
	Order     int
	Color     string
	CreatedAt time.Time
	UpdatedAt int32
}

func (m *HistoryProcess) TableName() string {
	return "history_processes"
}
