package facility_work_schedules

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"time"
)

func PostFacilityWorkSchedulesIdInvoke(c *gin.Context) (openapi_models.PostFacilityWorkSchedulesIdResponse, error) {

	facilityWorkScheduleRep := repository.NewFacilityWorkScheduleRepository(middleware.GetRepositoryMode(c)...)

	var facilityWorkScheduleReq openapi_models.PostFacilityWorkSchedulesRequest
	if err := c.ShouldBindJSON(&facilityWorkScheduleReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	facilityWorkScheduleRep.Upsert(db.FacilityWorkSchedule{
		Id:         facilityWorkScheduleReq.FacilityWorkSchedule.Id,
		FacilityId: facilityWorkScheduleReq.FacilityWorkSchedule.FacilityId,
		Date:       facilityWorkScheduleReq.FacilityWorkSchedule.Date,
		Type:       facilityWorkScheduleReq.FacilityWorkSchedule.Type,
		CreatedAt:  time.Time{},
		UpdatedAt:  0,
	})

	return openapi_models.PostFacilityWorkSchedulesIdResponse{}, nil

}
