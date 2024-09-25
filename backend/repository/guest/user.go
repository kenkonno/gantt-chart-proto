package guest

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"time"
)

// Auto generated start
func NewUserRepository() interfaces.UserRepositoryIF {
	return &userRepository{}
}

var guestUserId = int32(-1) // TODO: 本当は定数を使うべきだけど依存的に使えない。というか本来この実態はinteractor側にあったほうがいいと思う
var guestUser = db.User{
	Id:               &guestUserId,
	DepartmentId:     -1,
	LimitOfOperation: 0,
	LastName:         "ゲスト",
	FirstName:        "",
	Password:         "",
	Email:            "",
	Role:             "guest", // TODO: これも本当はだめ
	CreatedAt:        time.Time{},
	UpdatedAt:        0,
}
type userRepository struct {
}

func (r *userRepository) FindAll() []db.User {
	return []db.User{guestUser}
}

func (r *userRepository) Find(id int32) db.User {
	return guestUser
}

func (r *userRepository) FindByAuth(email string, password string) db.User {
	return guestUser
}
func (r *userRepository) FindByEmail(email string) db.User {
	return guestUser
}

func (r *userRepository) Upsert(m db.User) {
	panic("this is guest repository. can't upsert.")
}

func (r *userRepository) Delete(id int32) {
	panic("this is guest repository. can't delete.")
}
