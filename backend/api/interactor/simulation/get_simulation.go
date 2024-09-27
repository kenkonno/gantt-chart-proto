package simulation

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/simulation"
)

// GetSimulationInvoke シミュレーション状況を返却する。
func GetSimulationInvoke(c *gin.Context) openapi_models.GetSimulationResponse {

	simulationRep := simulation.NewSimulationLock()
	simulationLock := simulationRep.Find(constants.SimulateTypeSchedule)

	return openapi_models.GetSimulationResponse{
		SimulationLock: openapi_models.SimulationLock{
			SimulationName: simulationLock.SimulationName,
			Status:         simulationLock.Status,
			LockedAt:       simulationLock.LockedAt,
			LockedBy:       simulationLock.LockedBy,
		},
	}
}
