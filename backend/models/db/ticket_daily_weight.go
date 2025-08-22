package db

import "time"

type TicketDailyWeight struct {
	TicketId int32     `gorm:"uniqueIndex:idx_ticket_daily_weight_ticket_date"`
	Date     time.Time `gorm:"uniqueIndex:idx_ticket_daily_weight_ticket_date"`
	WorkHour  int32
	CreatedAt time.Time
	UpdatedAt int32
}
