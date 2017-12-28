package main

import (
	"github.com/caarlos0/env"
	"github.com/kudrykv/webhookproxy/app/handler"
	"github.com/kudrykv/webhookproxy/app/internal/log"
	"github.com/kudrykv/webhookproxy/app/types"
	"goji.io"
	"goji.io/pat"
	"net/http"
)

func main() {
	cfg := types.Server{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	mux := goji.NewMux()
	wswh := handler.NewHandler()

	mux.HandleFunc(pat.New("/websocket/:channel"), wswh.WebSocket)
	mux.HandleFunc(pat.Post("/webhook/:channel"), wswh.Webhook)

	log.Info("app is about to start")
	http.ListenAndServe(":"+cfg.Port, mux)
}
