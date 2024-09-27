package simulation

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/simulation"
	"time"
)

// PostSimulationInvoke シミュレーションを開始する。
func PostSimulationInvoke(c *gin.Context) openapi_models.PostSimulationResponse {

	simulationLockRep := simulation.NewSimulationLock()

	userId := middleware.GetUserId(c)

	// 排他ロックを獲得する 既に存在する場合はDB側でエラーとなる
	simulationLockRep.Upsert(
		db.SimulationLock{
			SimulationName: constants.SimulateTypeSchedule,
			Status:         constants.SimulateStatusInProgress,
			LockedAt:       time.Now(),
			LockedBy:       *userId,
		})

	// 全てのシミュレーションテーブルを初期化しデータをコピーする
	simulationRep := simulation.NewSimulationRepository()
	simulationRep.InitAllData()

	return openapi_models.PostSimulationResponse{}
}
