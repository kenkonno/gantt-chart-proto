package gantt_groups

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func DeleteGanttGroupsIdInvoke(c *gin.Context) openapi_models.DeleteGanttGroupsIdResponse {

	ganttGroupRep := repository.NewGanttGroupRepository()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	ganttGroupRep.Delete(int32(id))

	return openapi_models.DeleteGanttGroupsIdResponse{}

}
