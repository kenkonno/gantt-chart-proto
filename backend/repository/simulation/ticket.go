package simulation

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewSimulationTicketRepository() interfaces.TicketRepositoryIF {
	return &ticketRepository{
		con:   connection.GetCon(),
		table: "simulation_tickets",
	}
}

type ticketRepository struct {
	con *gorm.DB
	table string
}

func (r *ticketRepository) FindAll() []db.Ticket {
	var tickets []db.Ticket

	result := r.con.Table(r.table).Order(`"order" ASC`).Find(&tickets)
	if result.Error != nil {
		panic(result.Error)
	}
	return tickets
}

func (r *ticketRepository) FindByFacilityType(facilityTypes []string, facilityStatus []string) []db.Ticket {
	var tickets []db.Ticket
	builder := r.con.Table(r.table).Order(`"order" ASC`).
		Joins("INNER JOIN simulation_gantt_groups ON simulation_gantt_groups.id = simulation_tickets.gantt_group_id").
		Joins("INNER JOIN simulation_facilities ON simulation_facilities.id = simulation_gantt_groups.facility_id")
	if len(facilityTypes) > 0 {
		builder.Where("simulation_facilities.type IN ?", facilityTypes)
	}
	if len(facilityStatus) > 0 {
		builder.Where("simulation_facilities.status IN ?", facilityStatus)
	}
	builder.Find(&tickets)
	if builder.Error != nil {
		panic(builder.Error)
	}
	return tickets
}

func (r *ticketRepository) Find(id int32) db.Ticket {
	var ticket db.Ticket

	result := r.con.Table(r.table).First(&ticket, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return ticket
}

func (r *ticketRepository) Upsert(m db.Ticket) (db.Ticket, error) {
	var org db.Ticket
	if m.Id != nil {
		result := r.con.Table(r.table).First(&org, m.Id)
		if result.Error != nil {
			panic(result.Error)
		}
		if org.Id != nil && org.UpdatedAt > m.UpdatedAt {
			return m, connection.NewConflictError()
		}
		r.con.Table(r.table).Omit("memo").Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}, {Name: "created_at"}},
			UpdateAll: true,
		}).Save(&m) // created_atが上書きされるので除外する。またSaveじゃないとUpdatedAtが更新後で取れない
	} else {
		r.con.Table(r.table).Omit("memo").Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&m) // 新規の場合はCreateにする。
	}
	return m, nil
}

func (r *ticketRepository) UpdateMemo(m db.Ticket) (db.Ticket, error) {
	var org db.Ticket
	result := r.con.Table(r.table).First(&org, m.Id)
	if result.Error != nil {
		panic(result.Error)
	}
	if org.UpdatedAt > m.UpdatedAt {
		return m, connection.NewConflictError()
	}

	r.con.Table(r.table).Model(&m).Update("memo", m.Memo)
	return m, nil
}

func (r *ticketRepository) Delete(id int32) {
	r.con.Table(r.table).Where("id = ?", id).Delete(db.Ticket{})
}

// Auto generated end
func (r *ticketRepository) FindByGanttGroupIds(ganttGroupIds []int32) []db.Ticket {
	var tickets []db.Ticket

	result := r.con.Table(r.table).Where("gantt_group_id IN ?", ganttGroupIds).Order("simulation_tickets.order ASC").Find(&tickets)
	if result.Error != nil {
		panic(result.Error)
	}
	return tickets
}

func (r *ticketRepository) FindByUserIds(userIds []int32, facilityStatus string) []db.Ticket {
	var tickets []db.Ticket

	result := r.con.Table(r.table).Distinct().
		Joins("JOIN simulation_ticket_users ON simulation_tickets.id = simulation_ticket_users.ticket_id").
		Joins("JOIN simulation_gantt_groups ON simulation_tickets.gantt_group_id = simulation_gantt_groups.id").
		Joins("JOIN simulation_facilities   ON simulation_facilities.id = simulation_gantt_groups.facility_id").
		Where("simulation_ticket_users.user_id IN ?", userIds).
		Where("simulation_facilities.status = ?", facilityStatus).
		Order("simulation_tickets.order ASC").
		Find(&tickets)

	if result.Error != nil {
		panic(result.Error)
	}
	return tickets
}
