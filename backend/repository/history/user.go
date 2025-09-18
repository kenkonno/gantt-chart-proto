package history

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
)

func NewUserRepository(historyId int32) interfaces.UserRepositoryIF {
	return &userRepository{
		historyId: historyId,
		con:       connection.GetCon(),
	}
}

type userRepository struct {
	historyId int32
	con       *gorm.DB
}

func (r *userRepository) FindAll() []db.User {
	var users []db.User
	result := r.con.Table("history_users").Where("history_id = ?", r.historyId).Find(&users)
	if result.Error != nil {
		panic(result.Error)
	}
	return users
}

func (r *userRepository) Find(id int32) db.User {
	var user db.User
	result := r.con.Table("history_users").Where("history_id = ? AND id = ?", r.historyId, id).First(&user)
	if result.Error != nil {
		// TODO: 存在しない場合もpanicになるので、呼び出し元でハンドリングする
		panic(result.Error)
	}
	return user
}

func (r *userRepository) FindByAuth(email string, password string) db.User {
	var user db.User
	result := r.con.Table("history_users").Where("history_id = ? AND email = ? AND password = ?", r.historyId, email, password).First(&user)
	if result.Error != nil {
		panic(result.Error)
	}
	return user
}

func (r *userRepository) FindByEmail(email string) db.User {
	var user db.User
	result := r.con.Table("history_users").Where("history_id = ? AND email = ?", r.historyId, email).First(&user)
	if result.Error != nil {
		panic(result.Error)
	}
	return user
}

func (r *userRepository) Upsert(m db.User) {
	// History is read-only
}

func (r *userRepository) Delete(id int32) {
	// History is read-only
}
