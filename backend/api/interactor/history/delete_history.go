package history

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"strconv"
)

type DeleteHistoryInteractor struct{}

func NewDeleteHistoryInteractor() DeleteHistoryInteractor {
	return DeleteHistoryInteractor{}
}

func (i *DeleteHistoryInteractor) Execute(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(err)
	}

	historyRep := repository.NewHistoryRepository()
	err = historyRep.Delete(int32(id))
	if err != nil {
		panic(err)
	}
}
