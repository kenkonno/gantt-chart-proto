package interfaces

import "github.com/kenkonno/gantt-chart-proto/backend/models/db"

type ProcessRepositoryIF interface {
	FindAll() []db.Process
	Find(id int32) db.Process
	Upsert(m db.Process)
	Delete(id int32)
}
