package db

import "time"

type SimulationProcess struct {
	Id    *int32 `gorm:"primaryKey;autoIncrement"`
	Name  string
	Order int
	Color string

	CreatedAt time.Time
	UpdatedAt int32
}
