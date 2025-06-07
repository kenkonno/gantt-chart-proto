package simulation

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewSimulationFacilityWorkScheduleRepository() interfaces.FacilityWorkScheduleRepositoryIF {
	return &facilityWorkScheduleRepository{
		con:   connection.GetCon(),
		table: "simulation_facility_work_schedules",
	}
}

type facilityWorkScheduleRepository struct {
	con   *gorm.DB
	table string
}

func (r *facilityWorkScheduleRepository) FindByFacilityId(facilityId int32) []db.FacilityWorkSchedule {
	var facilityWorkSchedules []db.FacilityWorkSchedule

	result := r.con.Table(r.table).Where("facility_id = ?", facilityId).Order("date DESC").Find(&facilityWorkSchedules)
	if result.Error != nil {
		panic(result.Error)
	}
	return facilityWorkSchedules
}

func (r *facilityWorkScheduleRepository) FindAll() []db.FacilityWorkSchedule {
	var facilityWorkSchedules []db.FacilityWorkSchedule

	result := r.con.Table(r.table).Order("date DESC").Find(&facilityWorkSchedules)
	if result.Error != nil {
		panic(result.Error)
	}
	return facilityWorkSchedules
}

func (r *facilityWorkScheduleRepository) Find(id int32) db.FacilityWorkSchedule {
	var facilityWorkSchedule db.FacilityWorkSchedule

	result := r.con.Table(r.table).First(&facilityWorkSchedule, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return facilityWorkSchedule
}

func (r *facilityWorkScheduleRepository) Upsert(m db.FacilityWorkSchedule) db.FacilityWorkSchedule {
	r.con.Table(r.table).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
	return m
}

func (r *facilityWorkScheduleRepository) Delete(id int32) {
	r.con.Table(r.table).Where("id = ?", id).Delete(db.FacilityWorkSchedule{})
}

// Auto generated end
