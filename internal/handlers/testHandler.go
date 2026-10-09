package handlers

import (
	"fmt"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello from root handler")
}
