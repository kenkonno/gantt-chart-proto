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

	simulationLockRep := simulation.NewSimulationLock()

	if req.Mode == "apply" {
		// 各種simulation_xxx テーブルをリネームして本チャンに切り替える。
		simulationRep := simulation.NewSimulationRepository()
		simulationRep.SwitchTable()
		simulationRep.ResetSequence()
		simulationLockRep.Delete(constants.SimulateTypeSchedule)
	} else {
		simulationLock := simulationLockRep.Find(constants.SimulateTypeSchedule)
		if req.Mode == "pending" {
			simulationLock.Status = constants.SimulateStatusPending
		} else if req.Mode == "resume" {
			simulationLock.Status = constants.SimulateStatusInProgress
		}
		simulationLockRep.Upsert(simulationLock)
	}

	return openapi_models.PutSimulationResponse{}

}
