package repository

import (
	"fmt"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewHolidayRepository() holidayRepository {
	return holidayRepository{con}
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

	result := r.con.Where("facility_id = ?", facilityId).Order("id DESC").Find(&holidays)
	if result.Error != nil {
		panic(result.Error)
	}
	return holidays
}

func (r *holidayRepository) InsertByFacilityId(facilityId int32) []db.Holiday {
	var results []db.Holiday

	r.con.Raw(fmt.Sprintf(`
	WITH date_master AS (SELECT date                                                  as date,
								CASE
									WHEN extract(dow FROM date) = 6 THEN '土曜日'
									WHEN extract(dow FROM date) = 0 THEN '日曜日' END as youbi
						 FROM generate_series((SELECT term_from FROM facilities WHERE id = %d),
											  (SELECT term_to FROM facilities WHERE id = %d), '1 days') as date
						 WHERE extract(dow FROM date) IN (6, 0))
	INSERT
	INTO holidays (name, date, created_at, facility_id)
	SELECT youbi, date, now(), %d
	FROM date_master
	`, facilityId, facilityId, facilityId)).Scan(&results)
	return results
}
