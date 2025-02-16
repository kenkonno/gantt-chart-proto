package gantt_groups

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func DeleteGanttGroupsIdInvoke(c *gin.Context) (openapi_models.DeleteGanttGroupsIdResponse, error) {

	ganttGroupRep := repository.NewGanttGroupRepository(middleware.GetRepositoryMode(c)...)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	ganttGroupRep.Delete(int32(id))

	return openapi_models.DeleteGanttGroupsIdResponse{}, nil

}
