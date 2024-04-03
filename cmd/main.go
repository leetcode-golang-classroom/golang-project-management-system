package main

import (
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/leetcode-golang-classroom/golang-project-management-system/api"
	"github.com/leetcode-golang-classroom/golang-project-management-system/config"
	"github.com/leetcode-golang-classroom/golang-project-management-system/internal/db"
	"github.com/leetcode-golang-classroom/golang-project-management-system/internal/storage"
)

func main() {
	cfg := mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	sqlStorage := db.NewMySQLStorage(cfg)
	db, err := sqlStorage.Init()
	if err != nil {
		log.Fatal(err)
	}
	store := storage.NewStore(db)
	api_server := api.NewAPIServer(":3000", store)
	api_server.Serve()
}
