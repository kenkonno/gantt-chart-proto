package simulation

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewSimulationTicketUserRepository() interfaces.TicketUserRepositoryIF {
	return &ticketUserRepository{
		con:   connection.GetCon(),
		table: "simulation_ticket_users",
	}
}

type ticketUserRepository struct {
	con *gorm.DB
	table string
}

func (r *ticketUserRepository) FindAll() []db.TicketUser {
	var ticketUsers []db.TicketUser

	result := r.con.Table(r.table).Order("id DESC").Find(&ticketUsers)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticketUsers
}

func (r *ticketUserRepository) Find(id int32) db.TicketUser {
	var ticketUser db.TicketUser

	result := r.con.Table(r.table).First(&ticketUser, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticketUser
}
func (r *ticketUserRepository) FindByTicketId(ticketId int32) []db.TicketUser {
	var ticketUsers []db.TicketUser

	result := r.con.Table(r.table).Where("ticket_id = ? ", ticketId).Find(&ticketUsers)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticketUsers
}

func (r *ticketUserRepository) Upsert(m db.TicketUser) db.TicketUser {
	r.con.Table(r.table).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
	return m
}
func (r *ticketUserRepository) UpsertWithCreatedAt(m db.TicketUser) db.TicketUser {
	createdAt := m.CreatedAt
	result := r.Upsert(m)
	r.con.Table(r.table).Model(&result).Update("CreatedAt", createdAt)
	return result
}

func (r *ticketUserRepository) Delete(id int32) {
	r.con.Table(r.table).Where("id = ?", id).Delete(db.TicketUser{})
}

func (r *ticketUserRepository) DeleteByTicketId(id int32) {
	r.con.Table(r.table).Where("ticket_id = ?", id).Delete(db.TicketUser{})
}

// Auto generated end
func (r *ticketUserRepository) FindAllByTicketIds(ticketIds []int32) []db.TicketUser {
	var ticketUsers []db.TicketUser

	result := r.con.Table(r.table).Where("ticket_id IN ?", ticketIds).Order(`id, "order" `).Find(&ticketUsers)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticketUsers
}