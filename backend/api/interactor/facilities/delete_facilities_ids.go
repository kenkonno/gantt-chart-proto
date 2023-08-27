package facilities

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func DeleteFacilitiesIdInvoke(c *gin.Context) openapi_models.DeleteFacilitiesIdResponse {

	facilityRep := repository.NewFacilityRepository()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	facilityRep.Delete(int32(id))

	return openapi_models.DeleteFacilitiesIdResponse{}

}
