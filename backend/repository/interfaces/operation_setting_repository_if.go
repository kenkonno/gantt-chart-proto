package interfaces

import "github.com/kenkonno/gantt-chart-proto/backend/models/db"

type OperationSettingRepositoryIF interface {
	FindAll() []db.OperationSetting
	Find(id int32) db.OperationSetting
	Upsert(m db.OperationSetting)
	Delete(id int32)
	FindByFacilityId(facilityId int32) []db.OperationSetting
}
