package units

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
	"strconv"
)

func DeleteUnitsIdInvoke(c *gin.Context) openapi_models.DeleteUnitsIdResponse {

	unitRep := repository.NewUnitRepository()
	ganttGroupsRep := repository.NewGanttGroupRepository()
	ticketRep := repository.NewTicketRepository()
	operationSettingRep := repository.NewOperationSettingRepository()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	unit := unitRep.Find(int32(id))
	ganttGroups := lo.Filter(ganttGroupsRep.FindByFacilityId(unit.FacilityId), func(item db.GanttGroup, index int) bool {
		return item.UnitId == int32(id)
	})
	operationSettings := operationSettingRep.FindByFacilityId(unit.FacilityId)
	allTickets := ticketRep.FindByGanttGroupIds(
		lo.Map(ganttGroups, func(item db.GanttGroup, index int) int32 {
			return *item.Id
		}))

	// unitの削除 OK
	unitRep.Delete(int32(id))

	// ganttGroupsの削除 OK
	ganttGroupsRep.DeleteByUnitId(int32(id))

	// 関連チケットの削除
	for _, item := range allTickets {
		ticketRep.Delete(*item.Id)
	}
	// 稼働設定。該当のユニットのみを削除
	for _, item := range operationSettings {
		if item.UnitId == int32(id) {
			operationSettingRep.Delete(*item.Id)
		}
	}

	return openapi_models.DeleteUnitsIdResponse{}

}
