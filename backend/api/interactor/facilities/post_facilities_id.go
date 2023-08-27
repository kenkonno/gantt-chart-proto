package facilities

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"time"
)

func PostFacilitiesIdInvoke(c *gin.Context) openapi_models.PostFacilitiesIdResponse {

	facilityRep := repository.NewFacilityRepository()

	var facilityReq openapi_models.PostFacilitiesRequest
	if err := c.ShouldBindJSON(&facilityReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	facilityRep.Upsert(db.Facility{
		Id:        facilityReq.Facility.Id,
		Name:      facilityReq.Facility.Name,
		TermFrom:  facilityReq.Facility.TermFrom,
		TermTo:    facilityReq.Facility.TermTo,
		CreatedAt: time.Time{},
		UpdatedAt: 0,
	})

	return openapi_models.PostFacilitiesIdResponse{}

}
