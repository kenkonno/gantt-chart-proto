package db

import (
	"github.com/lib/pq"
	"time"
)

// TODO: tableは存在しない。Entityのほうが正しい。

type DefaultPileUp struct {
	Id                  *int32 `gorm:"primaryKey;autoIncrement"`
	FacilityId          int32
	DepartmentId        *int32
	Estimate            *int32
	NumberOfWorker      *int32
	StartDate           *time.Time
	EndDate             *time.Time
	UserIds             pq.Int32Array `gorm:"type:int32[]"`
	NumberOfWorkerByDay float32
	NumberOfWorkDay     int32
	ValidIndexes        pq.Int32Array `gorm:"type:int32[]"`
}
