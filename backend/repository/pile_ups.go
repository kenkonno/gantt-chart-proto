package repository

import (
	"github.com/kenkonno/gantt-chart-proto/backend/repository/common"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/history"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/simulation"
	"strconv"
)

func NewPileUpsRepository(mode ...string) interfaces.PileUpsRepositoryIF {
	if len(mode) >= 1 {
		if mode[0] == HistoryMode {
			historyId, err := strconv.Atoi(mode[1])
			if err != nil {
				panic(err)
			}
			return history.NewPileUpsRepository(int32(historyId))
		}
		if mode[0] == SimulationMode {
			return simulation.NewSimulationPileUpsRepository()
		}
	}
	return common.NewPileUpsRepository()
}

