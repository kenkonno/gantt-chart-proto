package facilities

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
)

func GetFacilitiesInvoke(c *gin.Context) openapi_models.GetFacilitiesResponse {
	facilityRep := repository.NewFacilityRepository()

	facilityList := facilityRep.FindAll()

	return openapi_models.GetFacilitiesResponse{
		List: lo.Map(facilityList, func(item db.Facility, index int) openapi_models.Facility {
			return openapi_models.Facility{
				Id:        item.Id,
				Name:      item.Name,
				TermFrom:  item.TermFrom,
				TermTo:    item.TermTo,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			}
		}),
	}
}
