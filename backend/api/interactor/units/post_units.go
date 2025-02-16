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

func PostUnitsInvoke(c *gin.Context) (openapi_models.PostUnitsResponse, error) {

	unitRep := repository.NewUnitRepository(middleware.GetRepositoryMode(c)...)
	ganttGroupsRep := repository.NewGanttGroupRepository(middleware.GetRepositoryMode(c)...)

	var unitReq openapi_models.PostUnitsRequest
	if err := c.ShouldBindJSON(&unitReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}
	r := unitRep.Upsert(db.Unit{
		Name:       strings.TrimSpace(unitReq.Unit.Name),
		FacilityId: unitReq.Unit.FacilityId,
		Order:      int(unitReq.Unit.Order),
		CreatedAt:  time.Time{},
		UpdatedAt:  0,
	})

	ganttGroupsRep.Upsert(db.GanttGroup{
		Id:         nil,
		FacilityId: unitReq.Unit.FacilityId,
		UnitId:     *r.Id,
		CreatedAt:  time.Time{},
		UpdatedAt:  0,
	})

	return openapi_models.PostUnitsResponse{}, nil

}
