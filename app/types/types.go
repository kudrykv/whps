package types

import (
	"encoding/json"
	"net/http"
	"time"
)

type Server struct {
	Port string `env:"PORT" envDefault:"8080"`
}

type Req struct {
	Id     string          `json:"id"`
	Time   time.Time       `json:"time"`
	Status int             `json:"status"`
	Header http.Header     `json:"headers"`
	Body   json.RawMessage `json:"body"`
}
