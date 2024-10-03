package db

import (
	"time"
)

type SimulationLock struct {
	SimulationName string    `gorm:"primaryKey"` // 未使用だが拡張のために設定
	Status         string    // pending, in_progress
	LockedAt       time.Time // The time at which the simulation was locked for execution
	LockedBy       int32     // The ID of the user/machine/instance that locked the simulation
}
