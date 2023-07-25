package repository

import (
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Auto generated start
func NewUserRepository() userRepository {
	return userRepository{con}
}

type userRepository struct {
	con *gorm.DB
}

func (r *userRepository) FindAll() []db.User {
	var users []db.User

	result := r.con.Order("id DESC").Find(&users)
	if result.Error != nil {
		panic(result.Error)
	}
	return users
}

func (r *userRepository) Find(id int32) db.User {
	var user db.User

	result := r.con.First(&user, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return user
}

func (r *userRepository) Upsert(m db.User) {
	r.con.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(&m)
}

func (r *userRepository) Delete(id int32) {
	r.con.Where("id = ?", id).Delete(db.User{})
}

// Auto generated end
