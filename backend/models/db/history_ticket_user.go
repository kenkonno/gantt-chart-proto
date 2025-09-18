package db

import "time"

type HistoryTicketUser struct {
	Id        *int32 `gorm:"primaryKey"`
	HistoryId int32  `gorm:"primaryKey"`
	TicketId  int32
	UserId    int32
	Order     int
	CreatedAt time.Time
	UpdatedAt int32
}

func (m *HistoryTicketUser) TableName() string {
	return "history_ticket_users"
}
