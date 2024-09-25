package simulation

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
	"strconv"
)

func DeleteSimulationInvoke(c *gin.Context) openapi_models.DeleteSimulationResponse {

	facilityRep := repository.NewFacilityRepository()
	ticketRep := repository.NewTicketRepository()
	ganttGroupRep := repository.NewGanttGroupRepository()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	facilityRep.Delete(int32(id))

	// 関連レコードの削除
	ganttGroups := ganttGroupRep.FindByFacilityId(int32(id))
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

	return openapi_models.DeleteSimulationResponse{}

}
