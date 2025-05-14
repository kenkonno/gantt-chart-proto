package common

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewFacilityWorkScheduleRepository() interfaces.FacilityWorkScheduleRepositoryIF {
	return &facilityWorkScheduleRepository{connection.GetCon()}
}

type facilityWorkScheduleRepository struct {
	con *gorm.DB
}

func (r *facilityWorkScheduleRepository) FindAll() []db.FacilityWorkSchedule {
	var facilityWorkSchedules []db.FacilityWorkSchedule

	result := r.con.Order("id DESC").Find(&facilityWorkSchedules)
	if result.Error != nil {
		panic(result.Error)
	}
	return facilityWorkSchedules
}

func (r *facilityWorkScheduleRepository) Find(id int32) db.FacilityWorkSchedule {
	var facilityWorkSchedule db.FacilityWorkSchedule

	result := r.con.First(&facilityWorkSchedule, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return facilityWorkSchedule
}

func (r *facilityWorkScheduleRepository) Upsert(m db.FacilityWorkSchedule) db.FacilityWorkSchedule {
	r.con.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
	return m
}

func (r *facilityWorkScheduleRepository) Delete(id int32) {
	r.con.Where("id = ?", id).Delete(db.FacilityWorkSchedule{})
}

// Auto generated end
