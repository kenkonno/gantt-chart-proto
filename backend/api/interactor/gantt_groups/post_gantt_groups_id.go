package gantt_groups

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"time"
)

func PostGanttGroupsIdInvoke(c *gin.Context) openapi_models.PostGanttGroupsIdResponse {

	ganttGroupRep := repository.NewGanttGroupRepository()

	var ganttGroupReq openapi_models.PostGanttGroupsRequest
	if err := c.ShouldBindJSON(&ganttGroupReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	ganttGroupRep.Upsert(db.GanttGroup{
		Id:         ganttGroupReq.GanttGroup.Id,
		FacilityId: ganttGroupReq.GanttGroup.FacilityId,
		UnitId:     ganttGroupReq.GanttGroup.UnitId,
		CreatedAt:  time.Time{},
		UpdatedAt:  0,
	})

	return openapi_models.PostGanttGroupsIdResponse{}

}
