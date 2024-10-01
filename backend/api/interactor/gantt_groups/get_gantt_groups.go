package gantt_groups

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
	"strconv"
)

func GetGanttGroupsInvoke(c *gin.Context) openapi_models.GetGanttGroupsResponse {

	mode := c.Query("mode")
	ganttGroupRep := repository.NewGanttGroupRepository(middleware.GetRepositoryMode(c)...)
	if mode == "prod" {
		ganttGroupRep = repository.NewGanttGroupRepository()
	}

	facilityId, err := strconv.Atoi(c.Query("facilityId"))
	if err != nil {
		panic(err)
	}

	ganttGroupList := ganttGroupRep.FindByFacilityId(int32(facilityId))

	return openapi_models.GetGanttGroupsResponse{
		List: lo.Map(ganttGroupList, func(item db.GanttGroup, index int) openapi_models.GanttGroup {
			return openapi_models.GanttGroup{
				Id:         item.Id,
				FacilityId: item.FacilityId,
				UnitId:     item.UnitId,
				CreatedAt:  item.CreatedAt,
				UpdatedAt:  item.UpdatedAt,
			}
		}),
	}
}
