package crontest

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func init() {
	http.HandleFunc("/crondoingitsjob", cronJob)
	http.HandleFunc("/", index)
}

func cronJob(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	log.Infof(ctx, "I'm doing my job. zip it shorty")
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello retards")
}
