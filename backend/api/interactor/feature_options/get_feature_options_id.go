package feature_options

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

func GetFeatureOptionsIdInvoke(c *gin.Context) (openapi_models.GetFeatureOptionsIdResponse, error) {
	featureOptionRep := repository.NewFeatureOptionRepository(middleware.GetRepositoryMode(c)...)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		panic(err)
	}

	featureOption := featureOptionRep.Find(int32(id))

	return openapi_models.GetFeatureOptionsIdResponse{
		FeatureOption: openapi_models.FeatureOption{
			Id:        featureOption.Id,
			Name:      featureOption.Name,
			Enabled:   featureOption.Enabled,
			CreatedAt: featureOption.CreatedAt,
			UpdatedAt: featureOption.UpdatedAt,
		},
	}, nil
}
