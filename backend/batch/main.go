package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kenkonno/gantt-chart-proto/backend/batch/update_stremaers"
)

func main() {

	// バッチはひな形だけ。何もいじっていないのでサンプルと思ってください。

	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}

	flag.Parse()

	args := flag.Args()

	fmt.Println("Batch Start [" + args[0] + "]")
	switch args[0] {
	case "UpdateStreamers":
		update_stremaers.Execute()
	}
	fmt.Println("Batch End   [" + args[0] + "]")

}

func getDate(args []string) *string {
	fmt.Println(args)
	var date *string
	if len(args) == 2 {
		date = &args[1]
	}
	return date
}
