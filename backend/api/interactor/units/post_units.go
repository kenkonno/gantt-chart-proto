package units

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"net/http"
	"time"
)

func PostUnitsInvoke(c *gin.Context) openapi_models.PostUnitsResponse {

	unitRep := repository.NewUnitRepository()

	var unitReq openapi_models.PostUnitsRequest
	if err := c.ShouldBindJSON(&unitReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}
	unitRep.Upsert(db.Unit{
		Name:       unitReq.Unit.Name,
		FacilityId: unitReq.Unit.FacilityId,
		CreatedAt:  time.Time{},
		UpdatedAt:  0,
	})

	return openapi_models.PostUnitsResponse{}

}
