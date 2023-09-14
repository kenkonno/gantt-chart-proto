package units

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
	"strconv"
)

func GetUnitsInvoke(c *gin.Context) openapi_models.GetUnitsResponse {
	unitRep := repository.NewUnitRepository()

	facilityId, err := strconv.Atoi(c.Query("facilityId"))
	if err != nil {
		panic(err)
	}
	unitList := unitRep.FindByFacilityId(int32(facilityId))

	return openapi_models.GetUnitsResponse{
		List: lo.Map(unitList, func(item db.Unit, index int) openapi_models.Unit {
			return openapi_models.Unit{
				Id:         item.Id,
				Name:       item.Name,
				FacilityId: item.FacilityId,
				Order:      int32(item.Order),
				CreatedAt:  item.CreatedAt,
				UpdatedAt:  item.UpdatedAt,
			}
		}),
	}
}
