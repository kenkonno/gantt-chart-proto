package holidays

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"strings"
	"time"
)

func PostHolidaysIdInvoke(c *gin.Context) openapi_models.PostHolidaysIdResponse {

	holidayRep := repository.NewHolidayRepository(middleware.GetRepositoryMode(c)...)

	var holidayReq openapi_models.PostHolidaysRequest
	if err := c.ShouldBindJSON(&holidayReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	holidayRep.Upsert(db.Holiday{
		Id:         holidayReq.Holiday.Id,
		FacilityId: holidayReq.Holiday.FacilityId,
		Name:       strings.TrimSpace(holidayReq.Holiday.Name),
		Date:       holidayReq.Holiday.Date,
		CreatedAt:  time.Time{},
		UpdatedAt:  0,
	})

	return openapi_models.PostHolidaysIdResponse{}

}
