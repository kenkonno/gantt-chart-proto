package simulation

import (
	"fmt"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewSimulationOperationSettingRepository() interfaces.OperationSettingRepositoryIF {
	return &operationSettingRepository{
		con:   connection.GetCon(),
		table: "simulation_operation_settings",
	}
}

type operationSettingRepository struct {
	con *gorm.DB
	table string
}

func (r *operationSettingRepository) FindAll() []db.OperationSetting {
	var operationSettings []db.OperationSetting

	result := r.con.Table(r.table).Order("id ASC").Find(&operationSettings)
	if result.Error != nil {
		panic(result.Error)
	}
	return operationSettings
}

func (r *operationSettingRepository) Find(id int32) db.OperationSetting {
	var operationSetting db.OperationSetting

	result := r.con.Table(r.table).First(&operationSetting, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return operationSetting
}

func (r *operationSettingRepository) Upsert(m db.OperationSetting) {
	r.con.Table(r.table).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "facility_id"}, {Name: "unit_id"}, {Name: "process_id"}},
		UpdateAll: true,
	}).Create(&m)
}

func (r *operationSettingRepository) Delete(id int32) {
	r.con.Table(r.table).Where("id = ?", id).Delete(db.OperationSetting{})
}

// Auto generated end
func (r *operationSettingRepository) FindByFacilityId(facilityId int32) []db.OperationSetting {
	var results []db.OperationSetting

	r.con.Table(r.table).Raw(fmt.Sprintf(`
	SELECT
		simulation_operation_settings.id
	,   %d facility_id
	,   simulation_units.id as unit_id
	,   simulation_processes.id as process_id
	,   COALESCE(simulation_operation_settings.work_hour, 8) as work_hour
	,   COALESCE(simulation_operation_settings.created_at, now()) as created_at
	,   simulation_operation_settings.updated_at
	FROM
		simulation_processes
		CROSS JOIN
		simulation_units
	LEFT JOIN
		simulation_operation_settings
	ON
		simulation_operation_settings.unit_id = simulation_units.id
	AND simulation_operation_settings.process_id = simulation_processes.id
	AND simulation_operation_settings.facility_id = %d
	WHERE
		simulation_units.facility_id = %d
	ORDER BY simulation_units.id, simulation_processes.id
	`, facilityId, facilityId, facilityId)).Scan(&results)
	return results
}
