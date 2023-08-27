package units

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
)

func GetUnitsInvoke(c *gin.Context) openapi_models.GetUnitsResponse {
	unitRep := repository.NewUnitRepository()

	unitList := unitRep.FindAll()

	return openapi_models.GetUnitsResponse{
		List: lo.Map(unitList, func(item db.Unit, index int) openapi_models.Unit {
			return openapi_models.Unit{
				Id:        item.Id,
				Name:      item.Name,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			}
		}),
	}
}
