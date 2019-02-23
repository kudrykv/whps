package main

import (
	"github.com/caarlos0/env"
	"github.com/kudrykv/whps/whps/handler"
	"github.com/kudrykv/whps/whps/internal/log"
	"github.com/kudrykv/whps/whps/types"
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

	log.Info("whps is about to start")
	http.ListenAndServe(":"+cfg.Port, mux)
}
