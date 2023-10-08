package processes

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func DeleteProcessesIdInvoke(c *gin.Context) openapi_models.DeleteProcessesIdResponse {

	processRep := repository.NewProcessRepository()
	ticketRep := repository.NewTicketRepository()
	operationSettingRep := repository.NewOperationSettingRepository()

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
			ticketRep.Upsert(item)
		}
	}
	// 稼働設定からは削除
	allOperationSettings := operationSettingRep.FindAll()
	for _, item := range allOperationSettings {
		if item.ProcessId == int32(id) {
			operationSettingRep.Delete(*item.Id)
		}
	}

	return openapi_models.DeleteProcessesIdResponse{}

}
