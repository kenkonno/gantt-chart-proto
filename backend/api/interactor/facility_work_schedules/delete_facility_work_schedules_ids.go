package facility_work_schedules

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func DeleteFacilityWorkSchedulesIdInvoke(c *gin.Context) (openapi_models.DeleteFacilityWorkSchedulesIdResponse, error) {

	facilityWorkScheduleRep := repository.NewFacilityWorkScheduleRepository()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	facilityWorkScheduleRep.Delete(int32(id))

	return openapi_models.DeleteFacilityWorkSchedulesIdResponse{}, nil

}
