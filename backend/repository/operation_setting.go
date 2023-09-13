package repository

import (
	"fmt"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewOperationSettingRepository() operationSettingRepository {
	return operationSettingRepository{con}
}

type operationSettingRepository struct {
	con *gorm.DB
}

func (r *operationSettingRepository) FindAll() []db.OperationSetting {
	var operationSettings []db.OperationSetting

	result := r.con.Order("id ASC").Find(&operationSettings)
	if result.Error != nil {
		panic(result.Error)
	}
	return operationSettings
}

func (r *operationSettingRepository) Find(id int32) db.OperationSetting {
	var operationSetting db.OperationSetting

	result := r.con.First(&operationSetting, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return operationSetting
}

func (r *operationSettingRepository) Upsert(m db.OperationSetting) {
	r.con.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "facility_id"}, {Name: "unit_id"}, {Name: "process_id"}},
		UpdateAll: true,
	}).Create(&m)
}

func (r *operationSettingRepository) Delete(id int32) {
	r.con.Where("id = ?", id).Delete(db.OperationSetting{})
}

// Auto generated end
func (r *operationSettingRepository) FindByFacilityId(facilityId int32) []db.OperationSetting {
	var results []db.OperationSetting

	r.con.Raw(fmt.Sprintf(`
	SELECT
		operation_settings.id
	,   %d facility_id
	,   units.id as unit_id
	,   processes.id as process_id
	,   COALESCE(operation_settings.work_hour, 8) as work_hour
	,   COALESCE(operation_settings.created_at, now()) as created_at
	,   operation_settings.updated_at
	FROM
		processes
		CROSS JOIN
		units
	LEFT JOIN
		operation_settings
	ON
		operation_settings.unit_id = units.id
	AND operation_settings.process_id = processes.id
	AND operation_settings.facility_id = %d
	WHERE
		units.facility_id = %d
	ORDER BY units.id, processes.id
	`, facilityId, facilityId, facilityId)).Scan(&results)
	return results
}
