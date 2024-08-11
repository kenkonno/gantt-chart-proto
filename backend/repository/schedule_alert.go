package repository

import (
	"github.com/kenkonno/gantt-chart-proto/backend/repository/common"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
)

func NewScheduleAlertRepository() interfaces.ScheduleAlertRepositoryIF {
	return common.NewScheduleAlertRepository()
}
