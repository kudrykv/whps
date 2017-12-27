package main

import (
	"goji.io"
	"net/http"
	"github.com/caarlos0/env"
	"github.com/kudrykv/webhookproxy/app/config"
	"goji.io/pat"
	"github.com/gorilla/websocket"
	"fmt"
	"sync"
	"io/ioutil"
	"encoding/json"
)

type Req struct {
	Header http.Header `json:"headers"`
	Body   []byte      `json:"body"`
}

func main()  {
	cfg := config.Server{}
	env.Parse(&cfg)

	sm := sync.Map{}

	mux := goji.NewMux()
	ws := websocket.Upgrader{}
	ws.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	mux.HandleFunc(pat.New("/websocket/:channel"), func(w http.ResponseWriter, r *http.Request) {
		ch := pat.Param(r, "channel")
		if _, ok := sm.Load(ch); ok {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		c, err := ws.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("err:", err)
			return
		}

		sm.Store(ch, c)

		defer sm.Delete(ch)
		defer c.Close()

		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				return
			}
		}
	})

	mux.HandleFunc(pat.Post("/webhook/:channel"), func(w http.ResponseWriter, r *http.Request) {
		ch := pat.Param(r, "channel")
		cTemp, ok := sm.Load(ch)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		c, ok := cTemp.(*websocket.Conn)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		bodyBytes, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		req := Req{Header: r.Header, Body: bodyBytes}
		bytes, err := json.Marshal(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = c.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(":" + cfg.Port, mux)
}
