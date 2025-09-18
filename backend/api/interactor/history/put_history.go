package history

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

type PutHistoryInteractor struct{}

func NewPutHistoryInteractor() PutHistoryInteractor {
	return PutHistoryInteractor{}
}

func (i *PutHistoryInteractor) Execute(c *gin.Context, req openapi_models.PutHistoriesIdRequest) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(err)
	}

	historyRep := repository.NewHistoryRepository()
	err = historyRep.UpdateName(int32(id), req.Name)
	if err != nil {
		panic(err)
	}
}
