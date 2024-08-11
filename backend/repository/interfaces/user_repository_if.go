package interfaces

import "github.com/kenkonno/gantt-chart-proto/backend/models/db"

type UserRepositoryIF interface {
	FindAll() []db.User
	Find(id int32) db.User
	FindByAuth(email string, password string) db.User
	FindByEmail(email string) db.User
	Upsert(m db.User)
	Delete(id int32)
}
