package common

import (
	"fmt"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// Auto generated start
func NewHolidayRepository() interfaces.HolidayRepositoryIF {
	return &holidayRepository{connection.GetCon()}
}

type holidayRepository struct {
	con *gorm.DB
}

func (r *holidayRepository) FindAll() []db.Holiday {
	var holidays []db.Holiday

	result := r.con.Order("id DESC").Find(&holidays)
	if result.Error != nil {
		panic(result.Error)
	}
	return holidays
}

func (r *holidayRepository) Find(id int32) db.Holiday {
	var holiday db.Holiday

	result := r.con.First(&holiday, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return holiday
}

func (r *holidayRepository) Upsert(m db.Holiday) {
	r.con.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
}

func (r *holidayRepository) Delete(id int32) {
	r.con.Where("id = ?", id).Delete(db.Holiday{})
}

// Auto generated end

func (r *holidayRepository) FindByFacilityId(facilityId int32) []db.Holiday {
	var holidays []db.Holiday

	result := r.con.Where("facility_id = ?", facilityId).Order("date DESC").Find(&holidays)
	if result.Error != nil {
		panic(result.Error)
	}
	return holidays
}

func (r *holidayRepository) InsertByFacilityId(facilityId int32, from *time.Time, to *time.Time) []db.Holiday {
	var results []db.Holiday

	var fromToWhere = ""

	if from != nil && to != nil {
		fromToWhere = fmt.Sprintf(`
	AND (date < '%s'::date OR '%s'::date < date)
	`, from.Format(time.RFC3339), to.Format(time.RFC3339))
	}

	r.con.Debug().Raw(fmt.Sprintf(`
	WITH date_master AS (SELECT date                                                  as date,
								CASE
									WHEN extract(dow FROM date) = 6 THEN '土曜日'
									WHEN extract(dow FROM date) = 0 THEN '日曜日' END as youbi
						 FROM generate_series((SELECT term_from FROM facilities WHERE id = %d),
											  (SELECT term_to FROM facilities WHERE id = %d), '1 days') as date
	                     WHERE extract(dow FROM date) IN (6, 0)
    	                 AND NOT EXISTS (SELECT * FROM holidays h WHERE h.facility_id = %d AND h.date = date.date)
						 %s
						)
	INSERT
	INTO holidays (name, date, created_at, facility_id, updated_at)
	SELECT youbi, date, now(), %d, EXTRACT(EPOCH FROM now())
	FROM date_master
	`, facilityId, facilityId, facilityId, fromToWhere, facilityId)).Scan(&results)
	return results
}
