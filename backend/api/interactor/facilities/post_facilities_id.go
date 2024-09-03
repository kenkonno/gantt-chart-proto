package facilities

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"strings"
	"time"
)

func PostFacilitiesIdInvoke(c *gin.Context) openapi_models.PostFacilitiesIdResponse {

	facilityRep := repository.NewFacilityRepository()
	holidayRep := repository.NewHolidayRepository()

	var facilityReq openapi_models.PostFacilitiesRequest
	if err := c.ShouldBindJSON(&facilityReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	facilityRep.Upsert(db.Facility{
		Id:              facilityReq.Facility.Id,
		Name:            strings.TrimSpace(facilityReq.Facility.Name),
		TermFrom:        facilityReq.Facility.TermFrom,
		TermTo:          facilityReq.Facility.TermTo,
		Order:           int(facilityReq.Facility.Order),
		Status:          facilityReq.Facility.Status,
		Type:            facilityReq.Facility.Type,
		ShipmentDueDate: facilityReq.Facility.ShipmentDueDate,
		CreatedAt:       time.Time{},
		UpdatedAt:       0,
	})

	return openapi_models.PostFacilitiesIdResponse{}

}
