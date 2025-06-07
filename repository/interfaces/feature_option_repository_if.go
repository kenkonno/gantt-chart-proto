package interfaces

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
)

type FeatureOptionRepositoryIF interface {
	FindAll() []db.FeatureOption
	Find(id int32) db.FeatureOption
	Upsert(m db.FeatureOption)
	Delete(id int32)
}
