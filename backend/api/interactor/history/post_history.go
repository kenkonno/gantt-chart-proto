package history

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
)

type PostHistoryInteractor struct{}

func NewPostHistoryInteractor() PostHistoryInteractor {
	return PostHistoryInteractor{}
}

// Execute PostHistoriesRequest と PostHistoriesResponse はopenapi.yamlに定義したあとで生成される
func (i *PostHistoryInteractor) Execute(c *gin.Context, req openapi_models.PostHistoriesRequest) openapi_models.PostHistoriesResponse {
	historyRep := repository.NewHistoryRepository()

	// The request model needs to be defined in OpenAPI spec.
	// Assuming it has FacilityId and Name.
	historyId, err := historyRep.CreateSnapshot(req.FacilityId, req.Name)
	if err != nil {
		// Proper error handling needed
		panic(err)
	}

	return openapi_models.PostHistoriesResponse{
		HistoryId: historyId,
	}
}
