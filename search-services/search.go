package search_service

import (
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", index)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "this is spart... I mean this is search service")
}
