package units

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func GetUnitsIdInvoke(c *gin.Context) openapi_models.GetUnitsIdResponse {
	unitRep := repository.NewUnitRepository()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	unit := unitRep.Find(int32(id))

	return openapi_models.GetUnitsIdResponse{
		Unit: openapi_models.Unit{
			Id:         unit.Id,
			Name:       unit.Name,
			FacilityId: unit.FacilityId,
			CreatedAt:  unit.CreatedAt,
			UpdatedAt:  unit.UpdatedAt,
		},
	}
}
