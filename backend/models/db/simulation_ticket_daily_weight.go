package db

import "time"

type SimulationTicketDailyWeight struct {
	TicketDailyWeight
	// 複合INDEXの場合名前を自前でつける必要があるので上書き
	TicketId int32     `gorm:"uniqueIndex:idx_sim_ticket_daily_weight_ticket_date"`
	Date     time.Time `gorm:"uniqueIndex:idx_sim_ticket_daily_weight_ticket_date"`
}
