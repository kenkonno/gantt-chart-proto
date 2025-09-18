package history

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

type GetHistoriesInteractor struct{}

func NewGetHistoriesInteractor() GetHistoriesInteractor {
	return GetHistoriesInteractor{}
}

func (i *GetHistoriesInteractor) Execute(c *gin.Context) openapi_models.GetHistoriesResponse {
	facilityIdStr := c.Query("facility_id")
	facilityId, err := strconv.Atoi(facilityIdStr)
	if err != nil {
		// handle error
		panic(err)
	}

	historyRep := repository.NewHistoryRepository()
	dbHistories := historyRep.FindByFacilityId(int32(facilityId))

	var apiHistories []openapi_models.History
	for _, dbHistory := range dbHistories {
		apiHistories = append(apiHistories, openapi_models.History{
			Id:         *dbHistory.Id,
			FacilityId: dbHistory.FacilityId,
			Name:       dbHistory.Name,
			CreatedAt:  dbHistory.CreatedAt,
			UpdatedAt:  dbHistory.UpdatedAt,
		})
	}

	return openapi_models.GetHistoriesResponse{
		Histories: apiHistories,
	}
}
