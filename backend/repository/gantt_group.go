package repository

import (
	"github.com/kenkonno/gantt-chart-proto/backend/repository/common"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
)

func NewGanttGroupRepository() interfaces.GanttGroupRepositoryIF {
	return common.NewGanttGroupRepository()
}

