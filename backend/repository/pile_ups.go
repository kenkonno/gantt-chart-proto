package repository

import (
	"github.com/kenkonno/gantt-chart-proto/backend/repository/common"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
)

func NewPileUpsRepository() interfaces.PileUpsRepositoryIF {
	return common.NewPileUpsRepository()
}

