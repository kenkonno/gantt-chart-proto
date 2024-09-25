package interfaces

import "github.com/kenkonno/gantt-chart-proto/backend/models/db"

type DepartmentRepositoryIF interface {
	FindAll() []db.Department
	Find(id int32) db.Department
	Upsert(m db.Department)
	Delete(id int32)
}
