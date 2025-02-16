package holidays

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func DeleteHolidaysIdInvoke(c *gin.Context) (openapi_models.DeleteHolidaysIdResponse, error) {

	holidayRep := repository.NewHolidayRepository(middleware.GetRepositoryMode(c)...)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	holidayRep.Delete(int32(id))

	return openapi_models.DeleteHolidaysIdResponse{}, nil

}
