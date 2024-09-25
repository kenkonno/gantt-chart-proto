package interfaces

import "github.com/kenkonno/gantt-chart-proto/backend/models/db"

type TicketUserRepositoryIF interface {
	FindAll() []db.TicketUser
	Find(id int32) db.TicketUser
	FindByTicketId(ticketId int32) []db.TicketUser
	Upsert(m db.TicketUser) db.TicketUser
	UpsertWithCreatedAt(m db.TicketUser) db.TicketUser
	Delete(id int32)
	DeleteByTicketId(id int32)
	FindAllByTicketIds(ticketIds []int32) []db.TicketUser
}
