package repository

import (
	"github.com/kenkonno/gantt-chart-proto/backend/repository/common"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/guest"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
)

const GuestMode = "guest"

func NewUserRepository(mode ...string) interfaces.UserRepositoryIF {
	if len(mode) >= 1 {
		if mode[0] == GuestMode {
			return guest.NewUserRepository()
		}
	}
	return common.NewUserRepository()
}
