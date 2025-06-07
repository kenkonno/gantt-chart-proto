package db

import "time"

type TicketDailyWeight struct {
	Id        *int32 `gorm:"primaryKey;autoIncrement"`
	TicketId  int32  `gorm:"index:idx_ticket_daily_weight_ticket_id"`
	WorkHour  int32
	Date      int32
	CreatedAt time.Time
	UpdatedAt int32
}
