package db

import "time"

type SimulationTicketUser struct {
	Id       *int32 `gorm:"primaryKey;autoIncrement"`
	TicketId int32
	UserId   int32

	Order     int
	CreatedAt time.Time
	UpdatedAt int32
}
