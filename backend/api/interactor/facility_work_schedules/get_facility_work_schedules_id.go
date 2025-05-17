package facility_work_schedules

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func GetFacilityWorkSchedulesIdInvoke(c *gin.Context) (openapi_models.GetFacilityWorkSchedulesIdResponse, error) {
	facilityWorkScheduleRep := repository.NewFacilityWorkScheduleRepository()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	facilityWorkSchedule := facilityWorkScheduleRep.Find(int32(id))

	return openapi_models.GetFacilityWorkSchedulesIdResponse{
		FacilityWorkSchedule: openapi_models.FacilityWorkSchedule{
			Id:         facilityWorkSchedule.Id,
			FacilityId: facilityWorkSchedule.FacilityId,
			Date:       facilityWorkSchedule.Date,
			Type:       facilityWorkSchedule.Type,
			CreatedAt:  facilityWorkSchedule.CreatedAt,
			UpdatedAt:  int(facilityWorkSchedule.UpdatedAt),
		},
	}, nil
}
