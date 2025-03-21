package processes

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

func PostProcessesInvoke(c *gin.Context) (openapi_models.PostProcessesResponse, error) {

	processRep := repository.NewProcessRepository(middleware.GetRepositoryMode(c)...)

	var processReq openapi_models.PostProcessesRequest
	if err := c.ShouldBindJSON(&processReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}
	processRep.Upsert(db.Process{
		Name:      strings.TrimSpace(processReq.Process.Name),
		Order:     int(processReq.Process.Order),
		Color:     processReq.Process.Color,
		CreatedAt: time.Time{},
		UpdatedAt: 0,
	})

	return openapi_models.PostProcessesResponse{}, nil

}
