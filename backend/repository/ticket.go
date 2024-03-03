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

	result := r.con.Order(`"order" ASC`).Find(&tickets)
	if result.Error != nil {
		panic(result.Error)
	}
	return tickets
}

func (r *ticketRepository) FindByFacilityType(facilityTypes []string, facilityStatus []string) []db.Ticket {
	var tickets []db.Ticket
	builder := r.con.Order(`"order" ASC`).
		Joins("INNER JOIN gantt_groups ON gantt_groups.id = tickets.gantt_group_id").
		Joins("INNER JOIN facilities ON facilities.id = gantt_groups.facility_id")
	if len(facilityTypes) > 0 {
		builder.Where("facilities.type IN ?", facilityTypes)
	}
	if len(facilityStatus) > 0 {
		builder.Where("facilities.status IN ?", facilityStatus)
	}
	builder.Debug().Find(&tickets)
	if builder.Error != nil {
		panic(builder.Error)
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

func (r *ticketRepository) Upsert(m db.Ticket) (db.Ticket, error) {
	var org db.Ticket
	if m.Id != nil {
		result := r.con.First(&org, m.Id)
		if result.Error != nil {
			panic(result.Error)
		}
		if org.Id != nil && org.UpdatedAt > m.UpdatedAt {
			return m, NewConflictError()
		}
	}

	r.con.Omit("memo").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Save(&m)
	return m, nil
}

func (r *ticketRepository) UpdateMemo(m db.Ticket) (db.Ticket, error) {
	var org db.Ticket
	result := r.con.First(&org, m.Id)
	if result.Error != nil {
		panic(result.Error)
	}
	if org.UpdatedAt > m.UpdatedAt {
		return m, NewConflictError()
	}

	r.con.Model(&m).Update("memo", m.Memo)
	return m, nil
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
