package db

import "time"

type HistoryDepartment struct {
	Id        *int32 `gorm:"primaryKey"`
	HistoryId int32  `gorm:"primaryKey"`
	Name      string
	Color     string `gorm:"default:'rgb(66, 165, 246)'"`
	Order     int
	CreatedAt time.Time
	UpdatedAt int32
}

// TODO: AIに一通り実装してもらった。model定義とクエリーが若干怪しいので見直してください。

func (m *HistoryDepartment) TableName() string {
	return "history_departments"
}
