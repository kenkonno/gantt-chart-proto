package db

import (
	"time"
)

type User struct {
	Id                  *int32 `gorm:"primaryKey;autoIncrement"`
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

	CreatedAt time.Time
	UpdatedAt int64
}
