package feature_options

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

func PostFeatureOptionsInvoke(c *gin.Context) (openapi_models.PostFeatureOptionsResponse, error) {

	featureOptionRep := repository.NewFeatureOptionRepository(middleware.GetRepositoryMode(c)...)

	var featureOptionReq openapi_models.PostFeatureOptionsRequest
	if err := c.ShouldBindJSON(&featureOptionReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}
	featureOptionRep.Upsert(db.FeatureOption{
		Name:      strings.TrimSpace(featureOptionReq.FeatureOption.Name),
		Enabled:   featureOptionReq.FeatureOption.Enabled,
		CreatedAt: time.Time{},
		UpdatedAt: 0,
	})

	return openapi_models.PostFeatureOptionsResponse{}, nil

}
