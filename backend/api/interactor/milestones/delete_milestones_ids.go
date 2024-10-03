package milestones
import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)
func DeleteMilestonesIdInvoke(c *gin.Context) openapi_models.DeleteMilestonesIdResponse {

	milestoneRep := repository.NewMilestoneRepository(middleware.GetRepositoryMode(c)...)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	milestoneRep.Delete(int32(id))

	return openapi_models.DeleteMilestonesIdResponse{}

}
