package handlers

import (
	"fmt"
	"log"
	"net/http"
)

type AuthHandler struct {
	l *log.Logger
}

func NewAuthHandler(l *log.Logger) *AuthHandler {
	return &AuthHandler{l}
}

func (h *AuthHandler) open(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	fmt.Fprintf(w, "Unauth API!\n")
}

func (h *AuthHandler) secret(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	fmt.Fprintf(w, "Auth API!\n")
}
