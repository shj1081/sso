package main

import (
	"log"
	"net/http"

	"github.com/shj1081/sso/internal/config"
	"github.com/shj1081/sso/internal/db"
	"github.com/shj1081/sso/internal/server"
	"github.com/shj1081/sso/internal/storer"
)

func main() {
	// 1) 환경 변수 및 Config 로드
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2) DB 연결
	database, err := db.NewDatabase(cfg.DBURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()
	log.Println("Database connection established")

	// 3) Storer 생성
	st := storer.NewMySQLStorer(database)

	// 4) Server 생성
	srv := server.NewServer(cfg, st)

	// 5) 서버 실행
	log.Printf("Starting server on %s\n", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, srv.RegisterRoutes()); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
