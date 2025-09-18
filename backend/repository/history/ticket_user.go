package history

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
)

func NewTicketUserRepository(historyId int32) interfaces.TicketUserRepositoryIF {
	return &ticketUserRepository{
		historyId: historyId,
		con:       connection.GetCon(),
	}
}

type ticketUserRepository struct {
	historyId int32
	con       *gorm.DB
}

func (r *ticketUserRepository) FindAll() []db.TicketUser {
	var ticketUsers []db.TicketUser
	result := r.con.Table("history_ticket_users").Where("history_id = ?", r.historyId).Order("id DESC").Find(&ticketUsers)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticketUsers
}

func (r *ticketUserRepository) Find(id int32) db.TicketUser {
	var ticketUser db.TicketUser
	result := r.con.Table("history_ticket_users").Where("history_id = ? AND id = ?", r.historyId, id).First(&ticketUser)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticketUser
}
func (r *ticketUserRepository) FindByTicketId(ticketId int32) []db.TicketUser {
	var ticketUsers []db.TicketUser
	result := r.con.Table("history_ticket_users").Where("history_id = ? AND ticket_id = ? ", r.historyId, ticketId).Find(&ticketUsers)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticketUsers
}

func (r *ticketUserRepository) Upsert(m db.TicketUser) db.TicketUser {
	// History is read-only
	return m
}
func (r *ticketUserRepository) UpsertWithCreatedAt(m db.TicketUser) db.TicketUser {
	// History is read-only
	return m
}

func (r *ticketUserRepository) Delete(id int32) {
	// History is read-only
}

func (r *ticketUserRepository) DeleteByTicketId(id int32) {
	// History is read-only
}

func (r *ticketUserRepository) FindAllByTicketIds(ticketIds []int32) []db.TicketUser {
	var ticketUsers []db.TicketUser
	result := r.con.Table("history_ticket_users").Where("history_id = ? AND ticket_id IN ?", r.historyId, ticketIds).Order(`id, "order" `).Find(&ticketUsers)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticketUsers
}
