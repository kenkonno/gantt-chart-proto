package simulation

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/simulation"
)

// PutSimulationInvoke モードにより切り替える。シミュレーションを再開する or シミュレーションを保留する
func PutSimulationInvoke(c *gin.Context) openapi_models.PutSimulationResponse {

	req := openapi_models.PutSimulationRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request body"})
		return openapi_models.PutSimulationResponse{}
	}

	simulationRep := simulation.NewSimulationLock()
	simulationLock := simulationRep.Find(constants.SimulateTypeSchedule)


	if req.Mode == "pending" {
		simulationLock.Status = constants.SimulateStatusPending
	}

	if req.Mode == "Resume" {
		simulationLock.Status = constants.SimulateStatusInProgress
	}

	simulationRep.Upsert(simulationLock)

	return openapi_models.PutSimulationResponse{}

}
