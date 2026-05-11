package main

import (
	"log"
	"net/http"

	"orlder-car-line/server/internal/infra/config"
	"orlder-car-line/server/internal/handler/linebot"
	"orlder-car-line/server/internal/handler/webapi/handler"
)

func main() {
	mux := http.NewServeMux()

	// LINE Bot
	mux.HandleFunc("POST /webhook", linebot.HandleWebhook)

	// リッチメニューをサーバー起動時に自動登録
	go linebot.SetupRichMenu()

	// Web API
	mux.HandleFunc("POST /api/login", handler.HandleLogin)
	mux.HandleFunc("POST /api/logout", handler.HandleLogout)
	mux.HandleFunc("GET /api/me", handler.HandleMe)

	addr := ":" + config.Port
	log.Printf("サーバー起動: %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}