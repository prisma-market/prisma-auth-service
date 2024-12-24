package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kihyun1998/prisma-market/prisma-auth-service/internal/config"
	"github.com/kihyun1998/prisma-market/prisma-auth-service/internal/handlers"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: $v", err)
	}

	// 라우터 설정
	r := mux.NewRouter()

	// 핸들러 설정
	authHandler := handlers.NewAuthHandler(cfg)
	// 라우트 등록
	r.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// 서버 시작
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
