package history

import (
	"fmt"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
)

func NewOperationSettingRepository(historyId int32) interfaces.OperationSettingRepositoryIF {
	return &operationSettingRepository{
		historyId: historyId,
		con:       connection.GetCon(),
	}
}

type operationSettingRepository struct {
	historyId int32
	con       *gorm.DB
}

func (r *operationSettingRepository) FindAll() []db.OperationSetting {
	var operationSettings []db.OperationSetting
	result := r.con.Table("history_operation_settings").Where("history_id = ?", r.historyId).Order("id ASC").Find(&operationSettings)
	if result.Error != nil {
		panic(result.Error)
	}
	return operationSettings
}

func (r *operationSettingRepository) Find(id int32) db.OperationSetting {
	var operationSetting db.OperationSetting
	result := r.con.Table("history_operation_settings").Where("history_id = ? AND id = ?", r.historyId, id).First(&operationSetting)
	if result.Error != nil {
		panic(result.Error)
	}
	return operationSetting
}

func (r *operationSettingRepository) Upsert(m db.OperationSetting) {
	// History is read-only
}

func (r *operationSettingRepository) Delete(id int32) {
	// History is read-only
}

func (r *operationSettingRepository) FindByFacilityId(facilityId int32) []db.OperationSetting {
	var results []db.OperationSetting

	r.con.Raw(fmt.Sprintf(`
	SELECT
		os.id
	,   %d facility_id
	,   u.id as unit_id
	,   p.id as process_id
	,   COALESCE(os.work_hour, 8) as work_hour
	,   COALESCE(os.created_at, now()) as created_at
	,   os.updated_at
	FROM
		history_processes p
		CROSS JOIN
		history_units u
	LEFT JOIN
		history_operation_settings os
	ON
		os.unit_id = u.id
	AND os.process_id = p.id
	AND os.facility_id = %d
	AND os.history_id = ?
	WHERE
		u.facility_id = %d
	AND u.history_id = ?
	AND p.history_id = ?
	ORDER BY u.id, p.id
	`, facilityId, facilityId, facilityId), r.historyId, r.historyId, r.historyId).Scan(&results)
	return results
}
