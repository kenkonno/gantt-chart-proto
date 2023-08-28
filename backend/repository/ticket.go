package repository

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewTicketRepository() ticketRepository {
	return ticketRepository{con}
}

type ticketRepository struct {
	con *gorm.DB
}

func (r *ticketRepository) FindAll() []db.Ticket {
	var tickets []db.Ticket

	result := r.con.Order("id DESC").Find(&tickets)
	if result.Error != nil {
		panic(result.Error)
	}
	return tickets
}

func (r *ticketRepository) Find(id int32) db.Ticket {
	var ticket db.Ticket

	result := r.con.First(&ticket, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticket
}

func (r *ticketRepository) Upsert(m db.Ticket) {
	r.con.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
}

func (r *ticketRepository) Delete(id int32) {
	r.con.Where("id = ?", id).Delete(db.Ticket{})
}

// Auto generated end
func (r *ticketRepository) FindByGanttGroupIds(ganttGroupIds []int32) []db.Ticket {
	var tickets []db.Ticket

	result := r.con.Where("gantt_group_id IN ?", ganttGroupIds).Order("tickets.order ASC").Find(&tickets)
	if result.Error != nil {
		panic(result.Error)
	}
	return tickets
}
