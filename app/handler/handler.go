package handler

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/kudrykv/webhookproxy/app/internal/httpshort"
	"github.com/kudrykv/webhookproxy/app/internal/log"
	"github.com/kudrykv/webhookproxy/app/internal/signals"
	"github.com/kudrykv/webhookproxy/app/types"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"goji.io/pat"
	"io/ioutil"
	"net/http"
	"reflect"
	"sync"
	"time"
)

type wswhHandler struct {
	sm sync.Map
	ws websocket.Upgrader
}

func NewHandler() *wswhHandler {
	ws := websocket.Upgrader{}
	ws.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	return &wswhHandler{
		sm: sync.Map{},
		ws: ws,
	}
}

func (h *wswhHandler) WebSocket(w http.ResponseWriter, r *http.Request) {
	ch := pat.Param(r, "channel")
	l := log.WithFields(logrus.Fields{
		"channel": ch,
		"wt":      "ws",
	})
	l.Info("connection received")

	if _, ok := h.sm.Load(ch); ok {
		l.Warn("double-used channel")
		httpshort.StringMessage(w, http.StatusForbidden, "channel is busy")
		return
	}

	c, err := h.ws.Upgrade(w, r, nil)
	if err != nil {
		l.WithField("err", err).Error("failed to upgrade connection")
		httpshort.StringMessage(w, http.StatusInternalServerError, "failed to upgrade connection")
		return
	}

	h.sm.Store(ch, c)
	l.WithField("time", time.Now()).Info("connection stored")

	defer h.sm.Delete(ch)
	defer c.Close()
	defer func() {
		l.WithField("time", time.Now()).Info("connection deleted")
	}()

	for {
		mt, messageBytes, err := c.ReadMessage()
		if err != nil {
			l.WithField("err", err).Error("failed to read the message")
			return
		}

		l = l.WithField("mt", mt)

		switch mt {
		case websocket.TextMessage:
			resp := types.Req{}
			if err := json.Unmarshal(messageBytes, &resp); err != nil {
				l.WithField("err", err).Error("failed to unmarshal response")
				continue
			}

			hachiko, ok := signals.Get(resp.Id)
			if !ok {
				l.WithField("id", resp.Id).Error("failed to get channel to response to")
				continue
			}

			l.WithField("id", resp.Id).Info("replay response to the service")
			hachiko <- &resp

		default:
			l.Warn(string(messageBytes))
		}
	}
}

func (h *wswhHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	ch := pat.Param(r, "channel")
	id := xid.New().String()
	l := log.WithFields(logrus.Fields{
		"channel": ch,
		"wt":      "wh",
		"id":      id,
	})

	cTemp, ok := h.sm.Load(ch)
	if !ok {
		httpshort.StringMessage(w, http.StatusNotFound, "no one listens on this channel")
		l.Warn("no websocket for the channel")
		return
	}

	c, ok := cTemp.(*websocket.Conn)
	if !ok {
		httpshort.StringMessage(w, http.StatusInternalServerError, "failed to cast type")
		l.WithField("type", reflect.TypeOf(c).String()).Error("failed to cast type")
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		httpshort.StringMessage(w, http.StatusInternalServerError, "failed to read body")
		l.WithField("err", err).Error("failed to read body")
		return
	}

	req := types.Req{
		Id:     id,
		Time:   time.Now(),
		Header: r.Header,
		Body:   bodyBytes,
	}

	hachiko, errSignalCreate := signals.Create(req.Id)
	if errSignalCreate != nil {
		l.WithField("err", err).Error("failed to create signal to wait for")
	}

	bytes, err := json.Marshal(req)
	if err != nil {
		httpshort.StringMessage(w, http.StatusInternalServerError, "failed to marshal response")
		l.WithField("err", err).Error("failed to marshal response")
		return
	}

	err = c.WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		httpshort.StringMessage(w, http.StatusInternalServerError, "failed to write message")
		l.WithField("err", err).Error("failed to write message")
		return
	}

	if errSignalCreate == nil {
		resp := <-hachiko
		if resp == nil {
			l.Error("never received a response")
			return
		}

		status := resp.Status
		if status == 0 {
			status = http.StatusOK
		}

		w.WriteHeader(status)
		w.Write(resp.Body)

		return
	}

	w.WriteHeader(http.StatusOK)
	l.Info("ok")
}
