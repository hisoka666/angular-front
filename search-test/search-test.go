package search_test

import (
	"fmt"
	"html/template"
	lg "log"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/appengine/search"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

func init() {
	http.HandleFunc("/", index)
	http.HandleFunc("/adduser", addUser)
	http.HandleFunc("/cariuser", cari)
}

type User struct {
	Nama      string
	Umur      int
	Pekerjaan string
}

type UserSearch struct {
	Nama      string
	Umur      float64
	Pekerjaan string
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index.html").ParseFiles("index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		lg.Fatalf("Error adalah: %v", err)
	}
}

func addUser(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	umur, _ := strconv.Atoi(r.FormValue("umur"))
	user := &User{
		Nama:      r.FormValue("user"),
		Umur:      umur,
		Pekerjaan: r.FormValue("job"),
	}

	key := datastore.NewIncompleteKey(ctx, "User", nil)
	if _, err := datastore.Put(ctx, key, user); err != nil {
		log.Errorf(ctx, "Gagal menambahkan database: %v", err)
		return
	}
	cari := &UserSearch{
		Nama:      r.FormValue("user"),
		Umur:      float64(umur),
		Pekerjaan: r.FormValue("job"),
	}
	index, err := search.Open("users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := index.Put(ctx, "", cari)
	if err != nil {
		log.Errorf(ctx, "Gagal menambahkan ke index: %v", err)
		return
	}
	fmt.Fprintln(w, id)
	fmt.Fprintln(w, "Berhasil menambahkan data")
	time.Sleep(time.Second * 2)
	http.Redirect(w, r, "/", 303)

}

func cari(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	index, err := search.Open("users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	namauser := r.FormValue("nama")
	que := "snippet('" + r.FormValue("nama") + "', Nama, 50)"
	log.Infof(ctx, "String nama adalah : %s", que)
	log.Infof(ctx, "Nama adalah: %v", namauser)
	// qu := "Nama: " + nama
	// fieldex := search.FieldExpression{
	// 	Name: "cariFieldNama",
	// 	Expr: que,
	// }
	// expre := &search.SearchOptions{
	// 	Limit:         0,
	// 	IDsOnly:       false,
	// 	Sort:          nil,
	// 	Fields:        []string{"Nama"},
	// 	Expressions:   []search.FieldExpression{fieldex},
	// 	Facets:        nil,
	// 	Refinements:   nil,
	// 	Cursor:        "",
	// 	Offset:        0,
	// 	CountAccuracy: 0,
	// }
	t := index.Search(ctx, namauser, nil)
	for {
		var us UserSearch
		_, err := t.Next(&us)
		if err == search.Done {
			break
		}
		if err != nil {
			log.Errorf(ctx, "Tidak bisa mencari: %v", err)
			break
		}

		fmt.Fprintf(w, "<li> %v </li>", us.Nama)
	}
}
