package facility_work_schedules

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
)

func GetFacilityWorkSchedulesInvoke(c *gin.Context) (openapi_models.GetFacilityWorkSchedulesResponse, error) {

	facilityWorkScheduleRep := repository.NewFacilityWorkScheduleRepository(middleware.GetRepositoryMode(c)...)
	facilityId, err := strconv.Atoi(c.Query("facilityId"))
	if err != nil {
		panic(err)
	}
	facilityWorkScheduleList := facilityWorkScheduleRep.FindByFacilityId(int32(facilityId))

	return openapi_models.GetFacilityWorkSchedulesResponse{
		List: lo.Map(facilityWorkScheduleList, func(item db.FacilityWorkSchedule, index int) openapi_models.FacilityWorkSchedule {
			return openapi_models.FacilityWorkSchedule{
				Id:         item.Id,
				FacilityId: item.FacilityId,
				Date:       item.Date,
				Type:       item.Type,
				CreatedAt:  item.CreatedAt,
				UpdatedAt: int32(item.UpdatedAt),
			}
		}),
	}, nil
}
