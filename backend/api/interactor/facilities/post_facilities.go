package facilities

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"time"
)

func PostFacilitiesInvoke(c *gin.Context) openapi_models.PostFacilitiesResponse {

	facilityRep := repository.NewFacilityRepository()
	holidayRep := repository.NewHolidayRepository()

	var facilityReq openapi_models.PostFacilitiesRequest
	if err := c.ShouldBindJSON(&facilityReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}
	newFacility := facilityRep.Upsert(db.Facility{
		Name:      facilityReq.Facility.Name,
		TermFrom:  facilityReq.Facility.TermFrom,
		TermTo:    facilityReq.Facility.TermTo,
		Order:     int(facilityReq.Facility.Order),
		CreatedAt: time.Time{},
		UpdatedAt: 0,
	})
	holidayRep.InsertByFacilityId(*newFacility.Id)

	return openapi_models.PostFacilitiesResponse{}

}
