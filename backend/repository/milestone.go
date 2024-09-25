package repository

import (
	"github.com/kenkonno/gantt-chart-proto/backend/repository/common"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
)

func NewMilestoneRepository() interfaces.MilestoneRepositoryIF {
	return common.NewMilestoneRepository()
}
