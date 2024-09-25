package simulation

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewSimulationUserRepository() interfaces.UserRepositoryIF {
	return &userRepository{
		con:   connection.GetCon(),
		table: "simulation_users",
	}
}

type userRepository struct {
	con *gorm.DB
	table string
}

func (r *userRepository) FindAll() []db.User {
	var users []db.User

	result := r.con.Table(r.table).Order("id ASC").Find(&users)
	if result.Error != nil {
		panic(result.Error)
	}
	return users
}

func (r *userRepository) Find(id int32) db.User {
	var user db.User

	result := r.con.Table(r.table).First(&user, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return user
}

func (r *userRepository) FindByAuth(email string, password string) db.User {
	var user db.User

	result := r.con.Table(r.table).Where("email = ? AND password = ?", email, password).Find(&user)
	if result.Error != nil {
		panic(result.Error)
	}
	return user
}
func (r *userRepository) FindByEmail(email string) db.User {
	var user db.User

	result := r.con.Table(r.table).Where("email = ?", email).Find(&user)
	if result.Error != nil {
		panic(result.Error)
	}
	return user
}

func (r *userRepository) Upsert(m db.User) {
	r.con.Table(r.table).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
}

func (r *userRepository) Delete(id int32) {
	r.con.Table(r.table).Where("id = ?", id).Delete(db.User{})
}

// Auto generated end
