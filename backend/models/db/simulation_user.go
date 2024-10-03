package db

import (
	"time"
)

type SimulationUser struct {
	Id               *int32 `gorm:"primaryKey;autoIncrement"`
	DepartmentId     int32
	LimitOfOperation float32
	LastName         string
	FirstName        string
	Password         string
	Email            string
	Role             string

	CreatedAt time.Time
	UpdatedAt int64
}
