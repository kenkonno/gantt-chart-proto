package db

import "time"

type FeatureOption struct {
	Id        *int32 `gorm:"primaryKey;autoIncrement"`
	Name      string
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt int32
}
