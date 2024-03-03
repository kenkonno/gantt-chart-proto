package repository

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewTicketUserRepository() ticketUserRepository {
	return ticketUserRepository{con}
}

type ticketUserRepository struct {
	con *gorm.DB
}

func (r *ticketUserRepository) FindAll() []db.TicketUser {
	var ticketUsers []db.TicketUser

	result := r.con.Order("id DESC").Find(&ticketUsers)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticketUsers
}

func (r *ticketUserRepository) Find(id int32) db.TicketUser {
	var ticketUser db.TicketUser

	result := r.con.First(&ticketUser, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticketUser
}
func (r *ticketUserRepository) FindByTicketId(ticketId int32) []db.TicketUser {
	var ticketUsers []db.TicketUser

	result := r.con.Where("ticket_id = ? ", ticketId).Find(&ticketUsers)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticketUsers
}

func (r *ticketUserRepository) Upsert(m db.TicketUser) db.TicketUser {
	r.con.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
	return m
}
func (r *ticketUserRepository) UpsertWithCreatedAt(m db.TicketUser) db.TicketUser {
	createdAt := m.CreatedAt
	result := r.Upsert(m)
	r.con.Model(&result).Update("CreatedAt", createdAt)
	return result
}

func (r *ticketUserRepository) Delete(id int32) {
	r.con.Where("id = ?", id).Delete(db.TicketUser{})
}

func (r *ticketUserRepository) DeleteByTicketId(id int32) {
	r.con.Where("ticket_id = ?", id).Delete(db.TicketUser{})
}

// Auto generated end
func (r *ticketUserRepository) FindAllByTicketIds(ticketIds []int32) []db.TicketUser {
	var ticketUsers []db.TicketUser

	result := r.con.Where("ticket_id IN ?", ticketIds).Order(`id, "order" `).Find(&ticketUsers)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticketUsers
}
