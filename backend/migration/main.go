package main

import (
	"fmt"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
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
}

func migrate[T any](model T) {
	fmt.Println("############# Migrate Start")
	err := con.AutoMigrate(model)
	if err != nil {
		panic(err)
	}
}
