package history

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
)

func NewTicketRepository(historyId int32) interfaces.TicketRepositoryIF {
	return &ticketRepository{
		historyId: historyId,
		con:       connection.GetCon(),
	}
}

type ticketRepository struct {
	historyId int32
	con       *gorm.DB
}

func (r *ticketRepository) FindAll() []db.Ticket {
	var tickets []db.Ticket
	result := r.con.Table("history_tickets").Where("history_id = ?", r.historyId).Order(`"order" ASC`).Find(&tickets)
	if result.Error != nil {
		panic(result.Error)
	}
	return tickets
}

func (r *ticketRepository) FindByFacilityType(facilityTypes []string, facilityStatus []string) []db.Ticket {
	var tickets []db.Ticket
	builder := r.con.Table("history_tickets").Where("history_id = ?", r.historyId).Order(`"order" ASC`).
		Joins("INNER JOIN history_gantt_groups ON history_gantt_groups.id = history_tickets.gantt_group_id AND history_gantt_groups.history_id = ?", r.historyId).
		Joins("INNER JOIN history_facilities ON history_facilities.id = history_gantt_groups.facility_id AND history_facilities.history_id = ?", r.historyId)
	if len(facilityTypes) > 0 {
		builder.Where("history_facilities.type IN ?", facilityTypes)
	}
	if len(facilityStatus) > 0 {
		builder.Where("history_facilities.status IN ?", facilityStatus)
	}
	builder.Find(&tickets)
	if builder.Error != nil {
		panic(builder.Error)
	}
	return tickets
}

func (r *ticketRepository) Find(id int32) db.Ticket {
	var ticket db.Ticket
	result := r.con.Table("history_tickets").Where("history_id = ? AND id = ?", r.historyId, id).First(&ticket)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticket
}

func (r *ticketRepository) Upsert(m db.Ticket) (db.Ticket, error) {
	// History is read-only
	return m, nil
}

func (r *ticketRepository) UpdateMemo(m db.Ticket) (db.Ticket, error) {
	// History is read-only
	return m, nil
}

func (r *ticketRepository) Delete(id int32) {
	// History is read-only
}

func (r *ticketRepository) FindByGanttGroupIds(ganttGroupIds []int32) []db.Ticket {
	var tickets []db.Ticket
	result := r.con.Table("history_tickets").Where("history_id = ? AND gantt_group_id IN ?", r.historyId, ganttGroupIds).Order("history_tickets.order ASC").Find(&tickets)
	if result.Error != nil {
		panic(result.Error)
	}
	return tickets
}

func (r *ticketRepository) FindByUserIds(userIds []int32, facilityStatus string) []db.Ticket {
	var tickets []db.Ticket
	result := r.con.Table("history_tickets").Distinct().
		Joins("JOIN history_ticket_users ON history_tickets.id = history_ticket_users.ticket_id AND history_ticket_users.history_id = ?", r.historyId).
		Joins("JOIN history_gantt_groups ON history_tickets.gantt_group_id = history_gantt_groups.id AND history_gantt_groups.history_id = ?", r.historyId).
		Joins("JOIN history_facilities ON history_facilities.id = history_gantt_groups.facility_id AND history_facilities.history_id = ?", r.historyId).
		Where("history_tickets.history_id = ? AND history_ticket_users.user_id IN ?", r.historyId, userIds).
		Where("history_facilities.status = ?", facilityStatus).
		Order("history_tickets.order ASC").
		Find(&tickets)

	if result.Error != nil {
		panic(result.Error)
	}
	return tickets
}
