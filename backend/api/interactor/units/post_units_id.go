package units

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

func PostUnitsIdInvoke(c *gin.Context) (openapi_models.PostUnitsIdResponse, error) {

	unitRep := repository.NewUnitRepository(middleware.GetRepositoryMode(c)...)

	var unitReq openapi_models.PostUnitsRequest
	if err := c.ShouldBindJSON(&unitReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	unitRep.Upsert(db.Unit{
		Id:         unitReq.Unit.Id,
		Name:       strings.TrimSpace(unitReq.Unit.Name),
		FacilityId: unitReq.Unit.FacilityId,
		Order:      int(unitReq.Unit.Order),
		CreatedAt:  time.Time{},
		UpdatedAt:  0,
	})

	return openapi_models.PostUnitsIdResponse{}, nil

}
