package departments

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func DeleteDepartmentsIdInvoke(c *gin.Context) openapi_models.DeleteDepartmentsIdResponse {

	departmentRep := repository.NewDepartmentRepository(middleware.GetRepositoryMode(c)...)
	ticketRep := repository.NewTicketRepository(middleware.GetRepositoryMode(c)...)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	departmentRep.Delete(int32(id))

	// 関連チケットの削除 FIXME: 全件操作なのでパフォーマンス問題がある
	allTickets := ticketRep.FindAll()
	for _, item := range allTickets {
		if item.DepartmentId != nil && *item.DepartmentId == int32(id) {
			item.DepartmentId = nil
			_, _ = ticketRep.Upsert(item) // 部署のNULL化なので新規無（厳密にはするべき）
		}
	}

	return openapi_models.DeleteDepartmentsIdResponse{}

}
