package interfaces

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"time"
)

type PileUpsRepositoryIF interface {
	GetDefaultPileUps(excludeFacilityId int32, facilityTypes []string) []db.DefaultPileUp
	// GetValidIndexUsers 全件とるのはよろしくないがスピード優先でいったんこうする
	GetValidIndexUsers(time.Time) []db.ValidIndexUser
}
