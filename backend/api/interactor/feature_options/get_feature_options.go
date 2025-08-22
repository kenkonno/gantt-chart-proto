package feature_options

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
)

func GetFeatureOptionsInvoke(c *gin.Context) (openapi_models.GetFeatureOptionsResponse, error) {
	featureOptionRep := repository.NewFeatureOptionRepository()

	featureOptionList := featureOptionRep.FindAll()

	return openapi_models.GetFeatureOptionsResponse{
		List: lo.Map(featureOptionList, func(item db.FeatureOption, index int) openapi_models.FeatureOption {
			return openapi_models.FeatureOption{
				Id:        item.Id,
				Name:      item.Name,
				Enabled:   item.Enabled,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			}
		}),
	}, nil
}
