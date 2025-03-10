package main

import (
	"log"

	"github.com/shj1081/sso/db"
	"github.com/shj1081/sso/sso/handler"
	"github.com/shj1081/sso/sso/server"
	"github.com/shj1081/sso/sso/storer"
)

func main() {

	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	defer db.Close()
	log.Println("database connection established")

	st := storer.NewMySQLStorer(db.GetDB())
	srv := server.NewServer(st)
	hdl := handler.NewHandler(srv)
	handler.RegisterRoutes(hdl)
	handler.StartServer(":8080")
}
