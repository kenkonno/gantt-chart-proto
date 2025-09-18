package history

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
)

func NewFacilityWorkScheduleRepository(historyId int32) interfaces.FacilityWorkScheduleRepositoryIF {
	return &facilityWorkScheduleRepository{
		historyId: historyId,
		con:       connection.GetCon(),
	}
}

type facilityWorkScheduleRepository struct {
	historyId int32
	con       *gorm.DB
}

func (r *facilityWorkScheduleRepository) FindAll() []db.FacilityWorkSchedule {
	var facilityWorkSchedules []db.FacilityWorkSchedule
	result := r.con.Table("history_facility_work_schedules").Where("history_id = ?", r.historyId).Order("date DESC").Find(&facilityWorkSchedules)
	if result.Error != nil {
		panic(result.Error)
	}
	return facilityWorkSchedules
}

func (r *facilityWorkScheduleRepository) FindByFacilityId(facilityId int32) []db.FacilityWorkSchedule {
	var facilityWorkSchedules []db.FacilityWorkSchedule
	result := r.con.Table("history_facility_work_schedules").Where("history_id = ? AND facility_id = ?", r.historyId, facilityId).Order("date DESC").Find(&facilityWorkSchedules)
	if result.Error != nil {
		panic(result.Error)
	}
	return facilityWorkSchedules
}

func (r *facilityWorkScheduleRepository) Find(id int32) db.FacilityWorkSchedule {
	var facilityWorkSchedule db.FacilityWorkSchedule
	result := r.con.Table("history_facility_work_schedules").Where("history_id = ? AND id = ?", r.historyId, id).First(&facilityWorkSchedule)
	if result.Error != nil {
		panic(result.Error)
	}
	return facilityWorkSchedule
}

func (r *facilityWorkScheduleRepository) Upsert(m db.FacilityWorkSchedule) db.FacilityWorkSchedule {
	// History is read-only
	return m
}

func (r *facilityWorkScheduleRepository) Delete(id int32) {
	// History is read-only
}
