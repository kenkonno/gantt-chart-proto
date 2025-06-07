package main

import (
	"fmt"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/connection"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var con = connection.GetConnection()

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
	migrate(db.SimulationLock{})
	migrate(db.FacilityWorkSchedule{})
	migrate(db.FeatureOption{})
	migrate(db.TicketDailyWeight{})

	// simulation
	migrate(db.SimulationDepartment{})
	migrate(db.SimulationFacility{})
	migrate(db.SimulationHoliday{})
	migrate(db.SimulationOperationSetting{})
	migrate(db.SimulationProcess{})
	migrate(db.SimulationUnit{})
	migrate(db.SimulationUser{})
	migrate(db.SimulationGanttGroup{})
	migrate(db.SimulationTicketUser{})
	migrate(db.SimulationTicket{})
	migrate(db.SimulationMilestone{})
	migrate(db.SimulationFacilitySharedLink{})
	migrate(db.SimulationFacilityWorkSchedule{})
	migrate(db.SimulationFeatureOption{})
	migrate(db.SimulationTicketDailyWeight{})

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
			DepartmentId:        0,
			LimitOfOperation:    0,
			LastName:            "システム",
			FirstName:           "管理者",
			Password:            string(hashedPassword),
			Email:               "admin",
			Role:                "admin",
			EmploymentStartDate: time.Time{},
			CreatedAt:           time.Time{},
			UpdatedAt:           0,
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
