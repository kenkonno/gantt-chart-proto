package repository

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/samber/lo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strings"
)

var con *gorm.DB

func init() {

	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	// connectionの取得
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo", host, user, password, dbname, port)
	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	con = d
}

func GetConnection() *gorm.DB {
	return con
}

func BeginTransaction() {
	con.Begin()
}

func Commit() {
	con.Commit()
}

func createInParam(arrStr []string) string {
	return strings.Join(lo.Map(arrStr, func(item string, index int) string {
		return fmt.Sprintf("'%s'", item)
	}), ",")
}
