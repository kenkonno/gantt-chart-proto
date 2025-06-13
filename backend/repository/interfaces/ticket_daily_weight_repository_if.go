package interfaces

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"time"
)

type TicketDailyWeightRepositoryIF interface {
	FindAll() []db.TicketDailyWeight
	Find(ticketId int32, date time.Time) db.TicketDailyWeight
	FindByTicketId(ticketId int32) []db.TicketDailyWeight
	Upsert(m db.TicketDailyWeight)
	Delete(ticketId int32, date time.Time)
}
