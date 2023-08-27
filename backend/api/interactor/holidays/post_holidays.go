package holidays

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"time"
)

func PostHolidaysInvoke(c *gin.Context) openapi_models.PostHolidaysResponse {

	holidayRep := repository.NewHolidayRepository()

	var holidayReq openapi_models.PostHolidaysRequest
	if err := c.ShouldBindJSON(&holidayReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}
	holidayRep.Upsert(db.Holiday{
		FacilityId: holidayReq.Holiday.FacilityId,
		Name:       holidayReq.Holiday.Name,
		Date:       holidayReq.Holiday.Date,
		CreatedAt:  time.Time{},
		UpdatedAt:  0,
	})

	return openapi_models.PostHolidaysResponse{}

}
