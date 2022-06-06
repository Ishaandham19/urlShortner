package handlers

import (
	"fmt"
	"net/http"
)

func secret(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprint(rw, "In authenticated api")
}

func open(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprint(rw, "In unauthenticated api")
}
