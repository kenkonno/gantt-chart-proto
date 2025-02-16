package processes

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func DeleteProcessesIdInvoke(c *gin.Context) (openapi_models.DeleteProcessesIdResponse, error) {

	processRep := repository.NewProcessRepository(middleware.GetRepositoryMode(c)...)
	ticketRep := repository.NewTicketRepository(middleware.GetRepositoryMode(c)...)
	operationSettingRep := repository.NewOperationSettingRepository(middleware.GetRepositoryMode(c)...)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	processRep.Delete(int32(id))

	// 関連チケットの削除 FIXME: 全件操作なのでパフォーマンス問題がある
	// チケットは対象の工程をnilにする
	allTickets := ticketRep.FindAll()
	for _, item := range allTickets {
		if item.ProcessId != nil && *item.ProcessId == int32(id) {
			item.ProcessId = nil
			_, _ = ticketRep.Upsert(item) // 工程をNULL化する処理なのでエラーハンドリングしない（厳密にはするべき）
		}
	}
	// 稼働設定からは削除
	allOperationSettings := operationSettingRep.FindAll()
	for _, item := range allOperationSettings {
		if item.ProcessId == int32(id) {
			operationSettingRep.Delete(*item.Id)
		}
	}

	return openapi_models.DeleteProcessesIdResponse{}, nil

}
