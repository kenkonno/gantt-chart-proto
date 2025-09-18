package repository

import (
	"github.com/kenkonno/gantt-chart-proto/backend/repository/common"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/history"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"strconv"
)

func NewFacilitySharedLinkRepository(mode ...string) interfaces.FacilitySharedLinkRepositoryIF {
	if len(mode) >= 1 {
		if mode[0] == HistoryMode {
			historyId, err := strconv.Atoi(mode[1])
			if err != nil {
				panic(err)
			}
			return history.NewFacilitySharedLinkRepository(int32(historyId))
		}
	}
	return common.NewFacilitySharedLinkRepository()
}

