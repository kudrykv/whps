package handler

import (
	"net/http"
	"os"
)

type verifyHandler struct {
}

func NewVerify() *verifyHandler {
	return &verifyHandler{}
}

func (h *verifyHandler) Verify(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`google-site-verification: ` + os.Getenv("GOOGLE_VERIFY")))
}
