package gantt_groups

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"time"
)

func PostGanttGroupsInvoke(c *gin.Context) openapi_models.PostGanttGroupsResponse {

	ganttGroupRep := repository.NewGanttGroupRepository()

	var ganttGroupReq openapi_models.PostGanttGroupsRequest
	if err := c.ShouldBindJSON(&ganttGroupReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}
	ganttGroupRep.Upsert(db.GanttGroup{
		FacilityId: ganttGroupReq.GanttGroup.FacilityId,
		UnitId:     ganttGroupReq.GanttGroup.UnitId,
		Order:      ganttGroupReq.GanttGroup.Order,
		CreatedAt:  time.Time{},
		UpdatedAt:  0,
	})

	return openapi_models.PostGanttGroupsResponse{}

}
