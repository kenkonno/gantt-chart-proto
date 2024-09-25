package repository

import (
	"github.com/kenkonno/gantt-chart-proto/backend/repository/common"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
)

func NewOperationSettingRepository() interfaces.OperationSettingRepositoryIF {
	return common.NewOperationSettingRepository()
}

