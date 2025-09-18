package history

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
)

func NewFacilityRepository(historyId int32) interfaces.FacilityRepositoryIF {
	return &facilityRepository{
		historyId: historyId,
		con:       connection.GetCon(),
		table:     "history_facilities",
	}
}

type facilityRepository struct {
	historyId int32
	con       *gorm.DB
	table     string
}

func (r *facilityRepository) FindAll(facilityTypes []string, facilityStatus []string, orderColumn string) []db.Facility {
	var facilities []db.Facility

	builder := r.con.Table(r.table).Where("history_id = ?", r.historyId).Order(`"` + orderColumn + `" ASC`)
	if len(facilityTypes) > 0 {
		builder.Where("type IN ?", facilityTypes)
	}
	if len(facilityStatus) > 0 {
		builder.Where("status IN ?", facilityStatus)
	}

	result := builder.Find(&facilities)
	if result.Error != nil {
		panic(result.Error)
	}
	return facilities
}

func (r *facilityRepository) Find(id int32) db.Facility {
	var facility db.Facility
	result := r.con.Table("history_facilities").Where("history_id = ? AND id = ?", r.historyId, id).First(&facility)
	if result.Error != nil {
		panic(result.Error)
	}
	return facility
}

func (r *facilityRepository) Upsert(m db.Facility) db.Facility {
	// History is read-only
	return m
}

func (r *facilityRepository) Delete(id int32) {
	// History is read-only
}
