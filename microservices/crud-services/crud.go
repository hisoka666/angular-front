package crud

import (
	"net/http"
	"fmt"
)

func init() {
	http.HandleFunc("/", index)
}

func index(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "This is crud-service")
}