package repository

import (
	"github.com/kenkonno/gantt-chart-proto/backend/repository/common"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/guest"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/history"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/simulation"
	"strconv"
)

func NewUserRepository(mode ...string) interfaces.UserRepositoryIF {
	if len(mode) >= 1 {
		if mode[0] == HistoryMode {
			historyId, err := strconv.Atoi(mode[1])
			if err != nil {
				// TODO: error handling
				panic(err)
			}
			return history.NewUserRepository(int32(historyId))
		}
		if mode[0] == GuestMode {
			return guest.NewUserRepository()
		}
		if mode[0] == SimulationMode {
			return simulation.NewSimulationUserRepository()
		}
	}
	return common.NewUserRepository()
}
