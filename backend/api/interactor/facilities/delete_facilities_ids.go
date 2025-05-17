package facilities

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
	"strconv"
)

func DeleteFacilitiesIdInvoke(c *gin.Context) (openapi_models.DeleteFacilitiesIdResponse, error) {

	facilityRep := repository.NewFacilityRepository(middleware.GetRepositoryMode(c)...)
	ticketRep := repository.NewTicketRepository(middleware.GetRepositoryMode(c)...)
	ganttGroupRep := repository.NewGanttGroupRepository(middleware.GetRepositoryMode(c)...)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	facilityRep.Delete(int32(id))

	// 関連レコードの削除
	ganttGroups := ganttGroupRep.FindByFacilityId([]int32{int32(id)})
	allTickets := ticketRep.FindByGanttGroupIds(
		lo.Map(ganttGroups, func(item db.GanttGroup, index int) int32 {
			return *item.Id
		}))
	for _, item := range ganttGroups {
		ganttGroupRep.Delete(*item.Id)
	}
	for _, item := range allTickets {
		ticketRep.Delete(*item.Id)
	}

	return openapi_models.DeleteFacilitiesIdResponse{}, nil

}
