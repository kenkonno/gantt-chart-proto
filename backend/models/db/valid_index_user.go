package db

import (
	"github.com/lib/pq"
	"time"
)

// TODO: tableは存在しない。Entityのほうが正しい。

type ValidIndexUser struct {
	Date       time.Time
	IsHoliday  bool
	UserIds    pq.Int32Array `gorm:"type:int32[]"`
	ValidIndex int32
}
