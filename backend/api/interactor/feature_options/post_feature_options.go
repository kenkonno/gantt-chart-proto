package feature_options

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
	"net/http"
	"strings"
	"time"
)

func PostFeatureOptionsInvoke(c *gin.Context) (openapi_models.PostFeatureOptionsResponse, error) {

	featureOptionRep := repository.NewFeatureOptionRepository()

	var featureOptionReq openapi_models.PostFeatureOptionsRequest
	if err := c.ShouldBindJSON(&featureOptionReq); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		panic(err)
	}

	// Nameが constants.feature_options に定義されている内容であるかチェック
	name := strings.TrimSpace(featureOptionReq.FeatureOption.Name)

	// 有効な Feature 値のリストを作成
	validFeatures := []string{
		string(constants.ScheduleSimulation),
		string(constants.UnitExpandCollapse),
		string(constants.UnitCopy),
		string(constants.MultiSelectFilter),
		string(constants.ProjectListFreeText),
		string(constants.ProjectListNameSort),
		string(constants.ProgressInput),
		string(constants.DelayNotification),
		string(constants.ResourceStackingView),
		string(constants.ResourceStackingGraph),
		string(constants.WorkloadWeighting),
	}

	// 名前が有効なリストに含まれているかチェック
	if !lo.Contains(validFeatures, name) {
		err := errors.New("invalid feature name: must be one of the predefined features in constants.feature_options")
		c.JSON(http.StatusBadRequest, err.Error())
		return openapi_models.PostFeatureOptionsResponse{}, err
	}

	featureOptionRep.Upsert(db.FeatureOption{
		Name:      name,
		Enabled:   featureOptionReq.FeatureOption.Enabled,
		CreatedAt: time.Time{},
		UpdatedAt: 0,
	})

	return openapi_models.PostFeatureOptionsResponse{}, nil

}
