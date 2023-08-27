package units

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func DeleteUnitsIdInvoke(c *gin.Context) openapi_models.DeleteUnitsIdResponse {

	unitRep := repository.NewUnitRepository()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	unitRep.Delete(int32(id))

	return openapi_models.DeleteUnitsIdResponse{}

}
