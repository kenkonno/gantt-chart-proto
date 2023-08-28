package gantt_groups

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func GetGanttGroupsIdInvoke(c *gin.Context) openapi_models.GetGanttGroupsIdResponse {
	ganttGroupRep := repository.NewGanttGroupRepository()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	ganttGroup := ganttGroupRep.Find(int32(id))

	return openapi_models.GetGanttGroupsIdResponse{
		GanttGroup: openapi_models.GanttGroup{
			Id:         ganttGroup.Id,
			FacilityId: ganttGroup.FacilityId,
			UnitId:     ganttGroup.UnitId,
			Order:      ganttGroup.Order,
			CreatedAt:  ganttGroup.CreatedAt,
			UpdatedAt:  ganttGroup.UpdatedAt,
		},
	}
}
