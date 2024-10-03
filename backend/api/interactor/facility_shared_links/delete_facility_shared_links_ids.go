package facility_shared_links
import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

// 未使用
func DeleteFacilitySharedLinksIdInvoke(c *gin.Context) openapi_models.DeleteFacilitySharedLinksIdResponse {

	facilitySharedLinkRep := repository.NewFacilitySharedLinkRepository(middleware.GetRepositoryMode(c)...)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	facilitySharedLinkRep.Delete(int32(id))

	return openapi_models.DeleteFacilitySharedLinksIdResponse{}

}
