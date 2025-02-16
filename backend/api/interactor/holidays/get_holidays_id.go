package holidays

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func GetHolidaysIdInvoke(c *gin.Context) (openapi_models.GetHolidaysIdResponse, error) {
	holidayRep := repository.NewHolidayRepository(middleware.GetRepositoryMode(c)...)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	holiday := holidayRep.Find(int32(id))

	return openapi_models.GetHolidaysIdResponse{
		Holiday: openapi_models.Holiday{
			Id:         holiday.Id,
			FacilityId: holiday.FacilityId,
			Name:       holiday.Name,
			Date:       holiday.Date,
			CreatedAt:  holiday.CreatedAt,
			UpdatedAt:  holiday.UpdatedAt,
		},
	}, nil
}
