package db

import "time"

type Ticket struct {
	Id              *int32 `gorm:"primaryKey;autoIncrement"`
	GanttGroupId    int32
	ProcessId       *int32
	DepartmentId    *int32
	LimitDate       time.Time
	Estimate        *int32
	DaysAfter       *int32
	StartDate       time.Time
	EndDate         time.Time
	ProgressPercent *int32
	Order           int32
	CreatedAt       time.Time
	UpdatedAt       int
}
