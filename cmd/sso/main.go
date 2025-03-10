package main

import (
	"log"

	"github.com/shj1081/sso/config"
	"github.com/shj1081/sso/db"
	"github.com/shj1081/sso/sso/handler"
	"github.com/shj1081/sso/sso/server"
	"github.com/shj1081/sso/sso/storer"
)

func main() {
	// 환경 변수 로드
	cfg := config.LoadConfig()

	// 데이터베이스 연결
	database, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	defer database.Close()
	log.Println("Database connection established")

	// 서버 및 핸들러 설정
	st := storer.NewMySQLStorer(database.GetDB())
	srv := server.NewServer(st)
	hdl := handler.NewHandler(srv)

	// 라우트 등록 및 서버 시작
	handler.RegisterRoutes(hdl)
	handler.StartServer(cfg.ServerAddress)
}
