package facilities

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func GetFacilitiesIdInvoke(c *gin.Context) openapi_models.GetFacilitiesIdResponse {
	facilityRep := repository.NewFacilityRepository(middleware.GetRepositoryMode(c)...)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	facility := facilityRep.Find(int32(id))

	return openapi_models.GetFacilitiesIdResponse{
		Facility: openapi_models.Facility{
			Id:        facility.Id,
			Name:      facility.Name,
			TermFrom:  facility.TermFrom,
			TermTo:    facility.TermTo,
			Order:     int32(facility.Order),
			CreatedAt: facility.CreatedAt,
			UpdatedAt: facility.UpdatedAt,
			Status:    facility.Status,
			Type:      facility.Type,
			ShipmentDueDate: facility.ShipmentDueDate,
		},
	}
}
