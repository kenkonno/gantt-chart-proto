package facilities

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

func PostFacilitiesIdInvoke(c *gin.Context) (openapi_models.PostFacilitiesIdResponse, error) {

	facilityRep := repository.NewFacilityRepository(middleware.GetRepositoryMode(c)...)

	var facilityReq openapi_models.PostFacilitiesRequest
	if err := c.ShouldBindJSON(&facilityReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	// No need to get old facility since we're not using it for holidays anymore

	facilityRep.Upsert(db.Facility{
		Id:              facilityReq.Facility.Id,
		Name:            strings.TrimSpace(facilityReq.Facility.Name),
		TermFrom:        facilityReq.Facility.TermFrom,
		TermTo:          facilityReq.Facility.TermTo,
		Order:           int(facilityReq.Facility.Order),
		Status:          facilityReq.Facility.Status,
		Type:            facilityReq.Facility.Type,
		FreeText: facilityReq.Facility.FreeText,
		ShipmentDueDate: facilityReq.Facility.ShipmentDueDate,
		CreatedAt:       time.Time{},
		UpdatedAt:       0,
	})

	// Holiday creation is now facility-independent

	return openapi_models.PostFacilitiesIdResponse{}, nil

}
