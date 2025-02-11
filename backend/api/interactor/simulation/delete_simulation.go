package simulation

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/simulation"
)

// DeleteSimulationInvoke シミュレーションをキャンセルする。
func DeleteSimulationInvoke(c *gin.Context) (openapi_models.DeleteSimulationResponse, error) {

	simulationRep := simulation.NewSimulationLock()
	simulationRep.Delete(constants.SimulateTypeSchedule)

	return openapi_models.DeleteSimulationResponse{}, nil

}
