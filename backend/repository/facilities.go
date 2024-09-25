package repository

import (
	"github.com/kenkonno/gantt-chart-proto/backend/repository/common"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/simulation"
)

func NewFacilityRepository(mode ...string) interfaces.FacilityRepositoryIF {
	if len(mode) >= 1 {
		if mode[0] == SimulationMode {
			return simulation.NewSimulationFacilityRepository()
		}
	}
	return common.NewFacilityRepository()
}

