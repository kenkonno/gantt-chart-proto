package db

import (
	"time"
)

type HistoryUser struct {
	Id                  *int32 `gorm:"primaryKey"`
	HistoryId           int32  `gorm:"primaryKey"`
	DepartmentId        int32
	LimitOfOperation    float32
	LastName            string
	FirstName           string
	Password            string
	Email               string
	Role                string
	PasswordReset       bool `gorm:"default:false"`
	EmploymentStartDate time.Time
	EmploymentEndDate   *time.Time
	CreatedAt           time.Time
	UpdatedAt           int64
}

func (m *HistoryUser) TableName() string {
	return "history_users"
}
