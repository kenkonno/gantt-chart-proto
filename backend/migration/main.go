package main

import (
	"fmt"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var con = repository.GetConnection()

func main() {
	migrate(db.Department{})
	migrate(db.Facility{})
	migrate(db.Holiday{})
	migrate(db.OperationSetting{})
	migrate(db.Process{})
	migrate(db.Unit{})
	migrate(db.User{})
	migrate(db.GanttGroup{})
	migrate(db.TicketUser{})
	migrate(db.Ticket{})
	migrate(db.Milestone{})
	migrate(db.FacilitySharedLink{})
	createDefaultUser()
}

func createDefaultUser() {
	userRep := repository.NewUserRepository()
	adminUser := userRep.FindByEmail("admin")
	if adminUser.Id == nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("defaultpassword"), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}

		userRep.Upsert(db.User{
			DepartmentId:     0,
			LimitOfOperation: 0,
			LastName:         "管理者",
			Password:         string(hashedPassword),
			Email:            "admin",
			Role:             "admin",
			CreatedAt:        time.Time{},
			UpdatedAt:        0,
		})
	}
}

func migrate[T any](model T) {
	fmt.Println("############# Migrate Start")
	err := con.AutoMigrate(model)
	if err != nil {
		panic(err)
	}
}
