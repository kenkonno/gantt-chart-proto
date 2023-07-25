package main

import (
	"fmt"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
)

var con = repository.GetConnection()

func main() {
	migrate(db.User{})
}

func migrate[T any](model T) {
	fmt.Println("############# Migrate Start")
	err := con.AutoMigrate(model)
	if err != nil {
		panic(err)
	}
}
