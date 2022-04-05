package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"google.golang.org/api/iterator"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
)

type tabel struct {
	BCP         []bcp    `json:"bcp"`
	TabelPasien []pasien `json:"tabpts"`
}
type bcp struct {
	Tanggal      time.Time `json:"tgl"`
	Kunjungan    []pasien  `json:"listpasien"`
	IKI1         int       `json:"iki1"`
	IKI2         int       `json:"iki2"`
	KegiatanLuar []pasien  `json:"kegiatanluar"`
	Shift        string    `json:"shift"`
}

type pasien struct {
	StatusServer string    `json:"stat"`
	TglKunjungan time.Time `json:"tgl"`
	ShiftJaga    string    `json:"shift"`
	ATS          string    `json:"ats"`
	Dept         string    `json:"dept"`
	NoCM         string    `json:"nocm"`
	NamaPasien   string    `json:"nama"`
	Diagnosis    string    `json:"diag"`
	IKI          string    `json:"iki"`
	LinkID       string    `json:"link"`
	TglAsli      time.Time `json:"tglasli"`
	TglLahir     time.Time `json:"tgllahir"`
}

// CatchDataJSON merupakan struct umum untuk
// "menangkap" json yang dikirim dari body
// request frontend
type catchDataJSON struct {
	Data1  string `json:"data01"`
	Data2  string `json:"data02"`
	Data3  string `json:"data03"`
	Data4  string `json:"data04"`
	Data5  string `json:"data05"`
	Data6  string `json:"data06"`
	Data7  string `json:"data07"`
	Data8  string `json:"data08"`
	Data9  string `json:"data09"`
	Data10 string `json:"data10"`
	Data11 string `json:"data11"`
}

// Staff adalah struct untuk Kind datastore
// Staff
type Staff struct {
	Email       string `json:"email"`
	NamaLengkap string `json:"nama"`
	LinkID      string `json:"link"`
	Peran       string `json:"peran"`
}

// DetailStaf adalah struct untuk Kind
// Datastore DetailStaf
type DetailStaf struct {
	NamaLengkap  string    `json:"nama"`
	NIP          string    `json:"nip"`
	NPP          string    `json:"npp"`
	GolonganPNS  string    `json:"golpns"`
	Alamat       string    `json:"alamat"`
	Bagian       string    `json:"bagian"`
	LinkID       string    `json:"link"`
	TanggalLahir time.Time `json:"tgl"`
	Umur         string    `json:"umur"`
	TargetIKI    float64   `json:"target"`
	Admin        bool      `json:"adm"`
	Kalender     string    `json:"kalender"`
}

type admDetail struct {
	Admin    DetailStaf `json:"detadm"`
	Member   []Staff    `json:"member"`
	Kalender string     `json:"kalender"`
}

type dataPasien struct {
	NamaPasien string    `json:"namapts"`
	NomorCM    string    `json:"nocm"`
	JenKel     string    `json:"jenkel"`
	Alamat     string    `json:"alamat"`
	TglDaftar  time.Time `json:"tgldaf"`
	TglLahir   time.Time `json:"tgllhr"`
	Umur       time.Time `json:"umur"`
	LinkID     string    `json:"link"`
}

type compareDataPts struct {
	OldData   dataPasien      `json:"old"`
	NewData   dataPasien      `json:"new"`
	Kunjungan kunjunganPasien `json:"kunjungan"`
}

func (c compareDataPts) isSame() bool {
	if c.NewData.NamaPasien == c.OldData.NamaPasien && c.NewData.TglLahir == c.OldData.TglLahir {
		return true
	}
	return false
}

type kunjunganPasien struct {
	Diagnosis     string    `json:"diag"`
	LinkID        string    `json:"link"`
	GolIKI        string    `json:"iki"`
	ATS           string    `json:"ats"`
	ShiftJaga     string    `json:"shift"`
	JamDatang     time.Time `json:"jam"`
	Dokter        string    `json:"dokter"`
	Hide          bool      `json:"hide"`
	JamDatangRiil time.Time `json:"jamriil"`
	Bagian        string    `json:"bagian"`
}

type kegiatanDokter struct {
	IDPasien        string    `json:"idpts"`
	NamaTindakan    string    `json:"tindakan"`
	NamaPasien      string    `json:"namapts"`
	Diagnosis       string    `json:"diag"`
	TglTindakan     time.Time `json:"tgltindakan"`
	KeyDataTindakan string    `json:"keytindakan"`
	Hide            bool      `json:"hide"`
}

type detailPasien struct {
	IDPasien      dataPasien        `json:"idpts"`
	ListKunjungan []kunjunganPasien `json:"kunjungan"`
}

func main() {
	http.Handle("/", middleFirst(http.HandlerFunc(indexHandler)))
	http.Handle("/login", middleFirst(http.HandlerFunc(loginHandler)))
	http.Handle("/home", middleFirst(http.HandlerFunc(homeHandler)))
	http.Handle("/profil", middleFirst(http.HandlerFunc(profilHandler)))
	http.Handle("/daftar-baru", middleFirst(http.HandlerFunc(daftarHandler)))
	http.Handle("/tambah-dokter", middleFirst(http.HandlerFunc(tambahDokter)))
	// http.Handle("/list-dokter", middleFirst(http.HandlerFunc(listDokter)))
	http.Handle("/get-info-nocm", middleFirst(http.HandlerFunc(getInfoNoCM)))
	http.Handle("/tambah-data-kunjungan", middleFirst(http.HandlerFunc(tambahDataKunjungan)))
	// http.Handle("/test-page", middleFirst(http.HandlerFunc(testPage)))
	// http.Handle("/test-search", middleFirst(http.HandlerFunc(testSearch)))
	http.Handle("/home-content", middleFirst(http.HandlerFunc(homeContentHandler)))
	http.Handle("/get-bcp", middleFirst(http.HandlerFunc(getBCP)))
	http.Handle("/get-kunjungan-pasien", middleFirst(http.HandlerFunc(getKunjunganPasien)))
	http.Handle("/edit-data-kunjungan", middleFirst(http.HandlerFunc(editDataKunjungan)))
	http.Handle("/hapus-data-kunjungan", middleFirst(http.HandlerFunc(hapusDataKunjungan)))
	http.Handle("/ubah-tanggal-kunjungan", middleFirst(http.HandlerFunc(ubahTanggalKunjungan)))
	http.Handle("/get-pdf-bcp", middleFirst(http.HandlerFunc(getPDF)))
	http.Handle("/get-detail-pasien", middleFirst(http.HandlerFunc(getDetailPasien)))
	http.Handle("/ubah-data-dokter", middleFirst(http.HandlerFunc(ubahDataDokter)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

/////////////////// Awal Functional Item ////////////////////////////

func middleFirst(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print(r.Header.Get("email") + " accessing " + r.URL.Path + " accessed at " + indonesiaNow().Format("15:04:05 Mon, 02/01/2006"))
		next.ServeHTTP(w, r)
	})
}

func indonesiaNow() time.Time {
	zone, _ := time.LoadLocation("Asia/Makassar")
	return time.Now().In(zone)
}

var days = [...]string{
	"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}

var months = [...]string{
	"Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

func genTemplate(n interface{}, temp string) (string, error) {
	b := new(bytes.Buffer)

	funcs := template.FuncMap{
		"jam": func(t time.Time) string {
			zone, _ := time.LoadLocation("Asia/Makassar")
			return t.In(zone).Format("15:04")
		},
		"inc": func(i int) int {
			return i + 1
		},
		"propercaps": func(input string) string {
			words := strings.Fields(input)
			smallwords, _ := readBucket("script-igdsanglah", "nama.txt")

			for index, word := range words {
				// strings.ToLower(word)
				if strings.Contains(smallwords, " "+word+" ") {
					words[index] = word
				} else {
					words[index] = strings.Title(strings.ToLower(word))
				}
			}
			return strings.Join(words, " ")
		},
		"properdiag": func(input string) string {
			words := strings.Fields(input)
			smallwords, _ := readBucket("script-igdsanglah", "diagnosis.txt")

			for index, word := range words {
				if strings.Contains(smallwords, " "+word+" ") {
					words[index] = word
				} else {
					words[index] = strings.Title(strings.ToLower(word))
				}
			}
			return strings.Join(words, " ")
		},
		// strtgl digunakan untuk membuat string tanggal dari
		// sebuah type Time dengan format 02/01/2006
		"strtgl": func(t time.Time) string {
			return t.In(zonaIndo()).Format("2006-01-02")
		},
		// istimezero digunakan untuk mencari tahu apakah type Time nol
		// digunakan untuk memberikan nilai false pada alur if dalam template
		"istimezero": func(t time.Time) bool {
			return t.IsZero()
		},
		"strjson": func(js interface{}) string {
			jso, err := json.Marshal(js)
			if err != nil {
				return "Gagal mengubah JSON"
			}
			return string(jso)
		},
		// strtglhari digunakan untuk membuat string tanggal disertai
		// nama hari dengan format Mon, 02/01/2006
		"strtglhari": func(t time.Time) string {
			return fmt.Sprintf("%s, %s", days[t.Weekday()][:3], t.Format("02/01/2006"))
		},
		// convstrjaga digunakan untuk mengubah string ShiftJaga yang berupa
		// angka menjadi Pagi, Sore dan Malam
		"convstrjaga": func(j string) string {
			var m string
			switch j {
			case "1":
				m = "Pagi"
			case "2":
				m = "Sore"
			case "3":
				m = "Malam"
			}
			return m
		},
		// tglbcp digunakan untuk membuat nama link bcp tiap bulan
		"tglbcp": func(tgl time.Time, shift string) string {
			if tgl.Hour() < 12 && shift == "3" {
				return tgl.AddDate(0, 0, -1).Format("2006/01")
			}
			return tgl.Format("2006/01")
		},
		"umur": func(lahir time.Time) string {
			skrng := timeNowIndonesia()
			yr, mn, dy := diffAge(lahir, skrng)
			return fmt.Sprintf("%d Tahun %d Bulan %d Hari", yr, mn, dy)
		},
		"gettgl": func(pts []pasien) string {
			return fmt.Sprintf("%s %d", months[pts[0].TglKunjungan.Month()-1], pts[0].TglKunjungan.Year())
			// return pts[0].TglKunjungan.Format("Jan, 2006")
		},
		"tgliki": func(tgl time.Time) string {
			return tgl.Format("02-01-2006")
		},
		"totaliki": func(bcp []bcp) string {
			var a, b int
			c := &a
			d := &b
			for _, v := range bcp {
				*c = *c + v.IKI1
				*d = *d + v.IKI2
			}

			total := float32(a)*0.0032 + float32(b)*0.01
			return fmt.Sprintf("%.4f", total)
		},
		"tgljscript": func(tgl time.Time) string {
			return tgl.Format("2006-01-02")
		},
		"angk": func(num float64) string {
			return fmt.Sprintf("%.4f", num)
		},
		"cutquote": func(str string) string {
			return str[1:]
		},
	}

	tmp := template.Must(template.New("template.html").Funcs(funcs).ParseFiles("template.html"))

	err := tmp.ExecuteTemplate(b, temp, n)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func errorCheck(w http.ResponseWriter, code int, msg string, err error) {
	w.WriteHeader(code)
	log.Print(indonesiaNow().Format("15:04:05 Mon, 02/01/2006"))
	log.Print(msg, err)
	fmt.Fprintf(w, msg, err)
}

func decodeReq(r *http.Request) catchDataJSON {
	defer r.Body.Close()
	js := catchDataJSON{}
	json.NewDecoder(r.Body).Decode(&js)
	return js
}

func decodeRes(r *http.Response) catchDataJSON {
	defer r.Body.Close()
	js := catchDataJSON{}
	json.NewDecoder(r.Body).Decode(&js)
	return js
}

func changeStringtoTime(tgl string) time.Time {
	str, _ := time.ParseInLocation("2006-1-02", tgl, zonaIndo())
	return str
}
func zonaIndo() *time.Location {
	zone, _ := time.LoadLocation("Asia/Makassar")
	return zone
}

func timeNowIndonesia() time.Time {
	zone, _ := time.LoadLocation("Asia/Makassar")
	now := time.Now()
	return now.In(zone)
}

//////////////////////Akhir functional item ////////////////////////////////////////

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	script, err := genTemplate(nil, "basic-template")
	if err != nil {
		errorCheck(w, 500, "Gagal Mengeksekusi Template, alasan : %s", err)
		return
	}

	w.WriteHeader(200)
	fmt.Fprint(w, script)
}

func diffAge(a, b time.Time) (year, month, day int) {
	if a.Location() != b.Location() {
		a = a.In(b.Location())
	}

	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)

	if day < 0 {
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}
	return
}

func readBucket(ember, data string) (string, error) {

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("GAgal membuat klien: %s", err)
		return "", err
		// TODO: handle error.
	}
	obj := client.Bucket(ember).Object(data)
	rc, err := obj.NewReader(ctx)
	if err != nil {
		log.Printf("Gagal membuat reader: %s", err)
		return "", err
	}
	slurp, err := ioutil.ReadAll(rc)
	rc.Close()
	if err != nil {
		log.Printf("Gagal membaca bucket: %s", err)
		return "", err
	}
	return string(slurp), nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.NotFound(w, r)
		return
	}
	script, err := genTemplate(nil, "login")
	if err != nil {
		errorCheck(w, 500, "Gagal Mengeksekusi Template, alasan : %s", err)
		return
	}

	w.WriteHeader(200)
	fmt.Fprint(w, script)

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/home" {
		http.NotFound(w, r)
		return
	}
	// cal, err := getCalBali()
	// if err != nil {
	// 	errorCheck(w, http.StatusInternalServerError, "Gagal membuat kalender: %v", err)
	// 	return
	// }
	// tang := map[string]string{
	// 	"tanggal": cal,
	// }
	// js, _ := json.Marshal(tang)
	u := r.URL.Query()
	var text string
	if u.Get("email") == "suryasedana@gmail.com" {
		var adm admDetail
		// adm.Kalender = string(js)
		script, err := genTemplate(adm, "home")
		if err != nil {
			errorCheck(w, 500, "Gagal Mengeksekusi Template, alasan : %s", err)
			return
		}

		text = script
	} else {
		det := DetailStaf{
			Admin: false,
		}
		// det.Kalender = string(js)
		script, err := genTemplate(det, "home")
		if err != nil {
			errorCheck(w, 500, "Gagal Mengeksekusi Template, alasan : %s", err)
			return
		}
		text = script
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, text)
}

func profilHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/profil" {
		http.NotFound(w, r)
		return
	}
	js := decodeReq(r)
	defer r.Body.Close()

	det, err := getDocProfile(js.Data1)
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal mengambil profil: %v", err)
		return
	}
	var resp string
	if js.Data1 == "suryasedana@gmail.com" {
		adm, err := getListStaf(det)
		if err != nil {
			errorCheck(w, 404, "Gagal mengambil daftar staf: %v", err)
			return
		}
		script, err := genTemplate(adm, "profil")
		if err != nil {
			errorCheck(w, 404, "Gagal membuat template: %v", err)
			return
		}
		resp = script
	} else {
		det.Admin = false
		script, err := genTemplate(det, "profil")
		if err != nil {
			errorCheck(w, 404, "Gagal membuat template: %v", err)
			return
		}
		resp = script
	}

	w.WriteHeader(200)
	fmt.Fprint(w, resp)
}

func getDocProfile(email string) (DetailStaf, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "igdsanglah")
	if err != nil {
		log.Printf("Could not create datastore client: %v", err)
	}
	query := datastore.NewQuery("Staff")
	it := client.Run(ctx, query)
	var det DetailStaf
	for {
		var st Staff
		ke, err := it.Next(&st)
		if err != nil && err == iterator.Done {
			break
		}
		if err != nil && err != iterator.Done {
			return det, err
		}
		if st.Email == email {
			q2 := datastore.NewQuery("DetailStaf").Ancestor(ke)
			ite := client.Run(ctx, q2)
			for {
				ku, err := ite.Next(&det)
				if err != nil && err == iterator.Done {
					break
				}
				det.LinkID = ku.Encode()
			}
			break
		}
	}
	return det, nil
}

func getListStaf(det DetailStaf) (admDetail, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "igdsanglah")
	if err != nil {
		log.Printf("Could not create datastore client: %v", err)
	}
	var list []Staff
	var adm admDetail
	query := datastore.NewQuery("Staff")
	_, err = client.GetAll(ctx, query, &list)
	if err != nil {
		return adm, err
	}
	adm.Admin = det
	adm.Member = list
	return adm, nil
}

func daftarHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/daftar-baru" {
		http.NotFound(w, r)
		return
	}
	script, err := genTemplate(nil, "daftar")
	if err != nil {
		errorCheck(w, 500, "Gagal Mengeksekusi Template, alasan : %s", err)
		return
	}

	w.WriteHeader(200)
	fmt.Fprint(w, script)
}

func tambahDokter(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/tambah-dokter" {
		http.NotFound(w, r)
		return
	}
	js := &catchDataJSON{}
	json.NewDecoder(r.Body).Decode(js)
	defer r.Body.Close()
	stf := Staff{
		Email:       js.Data1,
		Peran:       js.Data2,
		NamaLengkap: js.Data3,
	}
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "igdsanglah")
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal membuat klien datastore: %v", err)
		return
	}
	q := datastore.NewQuery("Staff").Filter("Email=", stf.Email)
	var staf []Staff
	_, err = client.GetAll(ctx, q, &staf)
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal mengecek ketersediaan dokter: %v", err)
		return
	}
	if len(staf) == 0 {
		k := datastore.NameKey("Staff", "", datastore.NameKey("IGD", "fasttrack", nil))
		_, err = client.Put(ctx, k, &stf)
		if err != nil {
			errorCheck(w, http.StatusInternalServerError, "Gagal menyimpan data: %v", err)
			return
		}
		w.WriteHeader(200)
	} else {
		errorCheck(w, http.StatusForbidden, "Alamat email sudah terdaftar, silahkan menggunakan alamat email yang lain", nil)
		return
	}
	fmt.Fprint(w, "Berhasil menyimpan data dokter baru")
}

func getInfoNoCM(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/get-info-nocm" {
		http.NotFound(w, r)
		return
	}
	js := decodeReq(r)
	r.Body.Close()

	pts := dataPasien{
		NomorCM: js.Data2,
	}
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "igdsanglah")
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal membuat klien datastore (CM): %v", err)
		return
	}
	var cpts compareDataPts
	k := datastore.NameKey("DataPasien", pts.NomorCM, datastore.NameKey("IGD", "fasttrack", nil))
	err = client.Get(ctx, k, &pts)
	if err != nil && err != datastore.ErrNoSuchEntity {
		errorCheck(w, http.StatusInternalServerError, "Gagal mengambil data pasien: kesalahan pada server %v", err)
		return
	} else if err == datastore.ErrNoSuchEntity {
		cp := compareDataPts{
			NewData: pts,
		}
		cpts = cp
	} else {
		pts.LinkID = k.Encode()
		cp := compareDataPts{
			OldData: pts,
			NewData: pts,
		}
		cpts = cp
	}
	script, err := genTemplate(cpts, "input-pts")
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal membuat template pasien: %v", err)
		return
	}

	w.WriteHeader(200)
	fmt.Fprint(w, script)
}

func tambahDataKunjungan(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/tambah-data-kunjungan" {
		http.NotFound(w, r)
		return
	}

	js := decodeReq(r)
	old := dataPasien{}
	err := json.Unmarshal([]byte(js.Data1), &old)
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal mengambil old data: %v", err)
		return
	}

	new := dataPasien{
		NamaPasien: js.Data3,
		TglLahir:   changeStringtoTime(js.Data4),
		TglDaftar:  timeNowIndonesia(),
		NomorCM:    js.Data2,
	}

	kun := kunjunganPasien{
		Diagnosis:     js.Data5,
		GolIKI:        js.Data8,
		ATS:           js.Data6,
		Bagian:        js.Data7,
		JamDatang:     timeNowIndonesia(),
		JamDatangRiil: timeNowIndonesia(),
		Dokter:        js.Data10,
		Hide:          false,
		ShiftJaga:     js.Data9,
	}
	comp := compareDataPts{
		NewData:   new,
		OldData:   old,
		Kunjungan: kun,
	}

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "igdsanglah")
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal membuat Client: %v", err)
		return
	}
	sam := comp.isSame()
	if comp.OldData.LinkID == "" {
		k := datastore.NameKey("DataPasien", comp.NewData.NomorCM, datastore.NameKey("IGD", "fasttrack", nil))
		_, err = client.Put(ctx, k, &comp.NewData)
		if err != nil {
			errorCheck(w, http.StatusInternalServerError, "Gagal menyimpan data pasien: %v", err)
			return
		}

		ka := datastore.NameKey("KunjunganPasien", "", k)
		_, err = client.Put(ctx, ka, &comp.Kunjungan)
		if err != nil {
			errorCheck(w, http.StatusInternalServerError, "Gagal menyimpan data kunjungan pasien: %v", err)
			return
		}
	} else if sam == false {
		k, err := datastore.DecodeKey(comp.OldData.LinkID)
		if err != nil {
			errorCheck(w, http.StatusInternalServerError, "Gagal mendecode kunci : %v", err)
			return
		}
		_, err = client.Put(ctx, k, &comp.NewData)
		if err != nil {
			errorCheck(w, http.StatusInternalServerError, "Gagal menyimpan data perubahan: %v", err)
			return
		}

		ka := datastore.NameKey("KunjunganPasien", "", k)
		_, err = client.Put(ctx, ka, &comp.Kunjungan)
		if err != nil {
			errorCheck(w, http.StatusInternalServerError, "Gagal menyimpan data kunjungan: %v", err)
			return
		}
	} else {
		k, err := datastore.DecodeKey(comp.OldData.LinkID)
		if err != nil {
			errorCheck(w, http.StatusInternalServerError, "Gagal mendecode kunci: %v", err)
			return
		}
		ke := datastore.NameKey("KunjunganPasien", "", k)
		_, err = client.Put(ctx, ke, &comp.Kunjungan)
		if err != nil {
			errorCheck(w, http.StatusInternalServerError, "Gagal menyimpan data kunjungan pasien: %v", err)
			return
		}
	}
	if js.Data11 != "" {
		keg := kegiatanDokter{
			NamaPasien:   js.Data3,
			NamaTindakan: js.Data11,
			Diagnosis:    js.Data5,
			TglTindakan:  kun.JamDatang,
			Hide:         false,
			IDPasien:     js.Data2,
		}

		q := datastore.NewQuery("Staff").Filter("Email=", kun.Dokter)
		var staf []Staff
		k, err := client.GetAll(ctx, q, &staf)
		if err != nil {
			errorCheck(w, http.StatusInternalServerError, "Nama dokter tidak ditemukan: %v", err)
			return
		}
		_, err = client.Put(ctx, datastore.NameKey("KegiatanDokter", "", k[0]), &keg)
		if err != nil {
			errorCheck(w, http.StatusInternalServerError, "gagal menyimpan data kegiatan dokter %v", err)
			return
		}
	}
	w.WriteHeader(200)
	fmt.Fprintf(w, "Berhasil menyimpan data kunjungan")
}

func homeContentHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/home-content" {
		http.NotFound(w, r)
		return
	}
	script, err := genTemplate(nil, "input-no-cm")
	if err != nil {
		errorCheck(w, 500, "Gagal Mengeksekusi Template, alasan : %s", err)
		return
	}
	w.WriteHeader(200)
	fmt.Fprint(w, script)
}

func getBCP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/get-bcp" {
		http.NotFound(w, r)
		return
	}

	js := decodeReq(r)
	r.Body.Close()

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "igdsanglah")
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal membuat klien: %v", err)
		return
	}

	bln, err := time.Parse("2006-1-02", js.Data2)
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal membuat tanggal awal: %v", err)
		return
	}
	q := datastore.NewQuery("KunjunganPasien").Filter("Dokter=", js.Data1).Filter("JamDatang <", bln.In(zonaIndo()).AddDate(0, 1, 0)).Filter("JamDatang >=", bln.In(zonaIndo())).Filter("Hide=", false).Order("-JamDatang")
	var kun []kunjunganPasien
	keys, err := client.GetAll(ctx, q, &kun)
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal mengambil data bcp: %v", err)
		return
	}

	if len(kun) == 0 {
		w.WriteHeader(404)
		script, err := genTemplate(nil, "pilih-bulan")
		if err != nil {
			errorCheck(w, http.StatusInternalServerError, "Gagal membuat template kunjungan kosong: %v", err)
			return
		}
		fmt.Fprint(w, script)
		return
	}
	list, err := iterateList(ctx, keys, kun, bln, client)
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal mengambil data pasien : %s", err)
		return
	}
	tab := getTabelPasien(list)

	script, err := genTemplate(tab, "tabel-bcp")
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal membuat tabel: %v", err)
		return
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, script)
}
func iterateList(ctx context.Context, keys []*datastore.Key, kun []kunjunganPasien, bln time.Time, cl *datastore.Client) ([]bcp, error) {
	var z []pasien
	for k, v := range kun {
		v.JamDatang = v.JamDatang.In(zonaIndo())
		if v.ShiftJaga == "3" && v.JamDatang.Hour() < 12 {
			v.JamDatang = v.JamDatang.AddDate(0, 0, -1)
		}
		if v.JamDatang.Month() != bln.Month() {
			continue
		}

		par := keys[k].Parent
		var n dataPasien
		err := cl.Get(ctx, par, &n)
		if err != nil {
			return nil, err
		}

		m := pasien{
			NamaPasien:   n.NamaPasien,
			TglKunjungan: v.JamDatang,
			ShiftJaga:    v.ShiftJaga,
			ATS:          v.ATS,
			Dept:         v.Bagian,
			NoCM:         par.String()[26:],
			Diagnosis:    v.Diagnosis,
			IKI:          v.GolIKI,
			LinkID:       keys[k].Encode(),
			TglAsli:      v.JamDatangRiil.In(zonaIndo()),
			TglLahir:     n.TglLahir.In(zonaIndo()),
		}

		z = append(z, m)
	}

	list := genListPasien(z)
	return list, nil
}

func genListPasien(pts []pasien) []bcp {
	g := []bcp{}
	for i := 1; i < 32; i++ {
		f := []pasien{}
		luar := []pasien{}
		h := bcp{}
		for _, b := range pts {
			if b.TglKunjungan.Day() == i {
				f = append(f, b)
				if b.NoCM == "00000000" || b.NoCM == "00000001" || b.NoCM == "00000002" || b.NoCM == "00000005" {
					luar = append(luar, b)
				}
			}

		}
		if len(f) != 0 {
			for i, j := 0, len(f)-1; i < j; i, j = i+1, j-1 {
				f[i], f[j] = f[j], f[i]
			}
			h.Tanggal = f[0].TglKunjungan
			h.Shift = f[(len(f) - 1)].ShiftJaga
		}
		h.Kunjungan = f
		h.KegiatanLuar = luar
		g = append(g, h)
	}

	for k, v := range g {
		var a, c int
		b := &a
		d := &c
		for _, n := range v.Kunjungan {
			if n.IKI == "1" {
				*b = *b + 1
			} else {
				*d = *d + 1
			}
		}

		g[k].IKI1 = a
		g[k].IKI2 = c
	}
	return g

}

func getTabelPasien(bcp []bcp) tabel {

	pts := []pasien{}
	for _, v := range bcp {
		for _, n := range v.Kunjungan {
			pts = append(pts, n)
		}
	}

	t := tabel{
		BCP:         bcp,
		TabelPasien: pts,
	}
	return t
}

func getKunjunganPasien(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/get-kunjungan-pasien" {
		http.NotFound(w, r)
		return
	}
	js := decodeReq(r)
	r.Body.Close()

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "igdsanglah")
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal membuat klien datastore", err)
	}
	// log.Printf("Key adalah: %s", js.Data2)
	k, err := datastore.DecodeKey(js.Data2)
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal mendecode key: %s", err)
	}

	kun := &kunjunganPasien{}

	if err := client.Get(ctx, k, kun); err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal mengambil data: %s", err)
	}
	zone, _ := time.LoadLocation("Asia/Makassar")
	kun.JamDatang = kun.JamDatang.In(zone)
	kun.JamDatangRiil = kun.JamDatangRiil.In(zone)
	kun.LinkID = k.Encode()
	script, err := genTemplate(kun, "edit-pasien")
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal membuat tabel: %v", err)
		return
	}
	w.WriteHeader(200)
	fmt.Fprintf(w, script)
}

func editDataKunjungan(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/edit-data-kunjungan" {
		http.NotFound(w, r)
		return
	}
	js := decodeReq(r)
	r.Body.Close()
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "igdsanglah")
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal membuat klien datastore", err)
	}
	// log.Printf("Key adalah: %s", js.Data2)
	k, err := datastore.DecodeKey(js.Data2)
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal mendecode key: %s", err)
	}
	// log.Printf("Key berhasil? %s", k)

	kun := &kunjunganPasien{}

	if err := client.Get(ctx, k, kun); err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal mengambil data: %s", err)
	}
	kun.LinkID = k.Encode()
	// log.Printf("Data adalah: %v", kun)
	kun.ATS = js.Data4
	kun.Bagian = js.Data5
	kun.Diagnosis = js.Data3
	kun.ShiftJaga = js.Data7
	kun.GolIKI = js.Data6
	zone, _ := time.LoadLocation("Asia/Makassar")
	kun.JamDatangRiil = kun.JamDatangRiil.In(zone)
	// log.Printf("Data baru adalah: %v", kun)
	_, err = client.Put(ctx, k, kun)
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "gagal menyimpan perubahan ke server: %s", err)
		return
	}
	w.WriteHeader(200)
}

func hapusDataKunjungan(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hapus-data-kunjungan" {
		http.NotFound(w, r)
		return
	}
	js := decodeReq(r)
	r.Body.Close()

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "igdsanglah")
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal membuat klien datastore", err)
		return
	}

	k, err := datastore.DecodeKey(js.Data2)
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal mendecode key: %s", err)
		return
	}
	kun := &kunjunganPasien{}

	if err := client.Get(ctx, k, kun); err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal mengambil data: %s", err)
		return
	}

	kun.Hide = true

	_, err = client.Put(ctx, k, kun)
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "gagal menyimpan perubahan ke server: %s", err)
		return
	}
	w.WriteHeader(200)
}

func ubahTanggalKunjungan(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ubah-tanggal-kunjungan" {
		http.NotFound(w, r)
		return
	}
	js := decodeReq(r)
	r.Body.Close()

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "igdsanglah")
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal membuat klien datastore", err)
		return
	}
	k, err := datastore.DecodeKey(js.Data2)
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal mendecode key: %s", err)
		return
	}

	kun := &kunjunganPasien{}

	if err := client.Get(ctx, k, kun); err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal mengambil data: %s", err)
		return
	}
	jam, err := time.Parse("2006-01-02", js.Data3)
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal memformat tanggal: %s", err)
		return
	}
	n := time.Now()
	zone, _ := time.LoadLocation("Asia/Makassar")
	switch js.Data4 {
	case "1":
		kun.JamDatang = time.Date(jam.Year(), jam.Month(), jam.Day(), 10, n.Minute(), n.Second(), n.Nanosecond(), zone)
	case "2":
		kun.JamDatang = time.Date(jam.Year(), jam.Month(), jam.Day(), 18, n.Minute(), n.Second(), n.Nanosecond(), zone)
	case "3":
		kun.JamDatang = time.Date(jam.Year(), jam.Month(), jam.Day(), 22, n.Minute(), n.Second(), n.Nanosecond(), zone)
	default:
		kun.JamDatang = time.Date(jam.Year(), jam.Month(), jam.Day(), 12, n.Minute(), n.Second(), n.Nanosecond(), zone)
	}
	kun.JamDatangRiil = kun.JamDatangRiil.In(zone)
	_, err = client.Put(ctx, k, kun)
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "gagal menyimpan perubahan ke server: %s", err)
		return
	}
	w.WriteHeader(200)
}

func getPDF(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/get-pdf-bcp" {
		http.NotFound(w, r)
		return
	}
	js := catchDataJSON{
		Data1: r.FormValue("email"),
		Data2: r.FormValue("tanggal"),
	}

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "igdsanglah")
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal membuat klien: %v", err)
		return
	}

	bln, err := time.Parse("2006-1-02", js.Data2)
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal membuat tanggal awal: %v", err)
		return
	}
	q := datastore.NewQuery("KunjunganPasien").Filter("Dokter=", js.Data1).Filter("JamDatang <", bln.In(zonaIndo()).AddDate(0, 1, 0)).Filter("JamDatang >=", bln.In(zonaIndo())).Filter("Hide=", false).Order("-JamDatang")
	var kun []kunjunganPasien
	keys, err := client.GetAll(ctx, q, &kun)
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal mengambil data bcp: %v", err)
		return
	}
	list, err := iterateList(ctx, keys, kun, bln, client)
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal mengambil data pasien : %s", err)
		return
	}
	det, err := getDocProfile(js.Data1)
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal mengambil profil: %v", err)
		return
	}
	t, err := createPDF(det, list)
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal membuat PDF: %s", err)
		return
	}
	w.Header().Set("Content-type", "application/pdf")
	if _, err := t.WriteTo(w); err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal membuat PDF: %s", err)
		return
	}
}

func createPDF(staf DetailStaf, list []bcp) (*bytes.Buffer, error) {
	shift := map[string]string{
		"1": "P",
		"2": "S",
		"3": "M",
	}
	a, b, c, d, e, f, g := getResumeIKI(list)
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 12)
	var tgl time.Time
	for _, v := range list {
		if v.Tanggal.IsZero() == false {
			tgl = v.Tanggal
			break
		}
	}
	// Tabel IKI \\\\\\\\\\\\\\\\///////////////////////////////////////////////
	pdf.AddPageFormat("L", gofpdf.SizeType{Wd: 210, Ht: 297})
	pdf.Cell(160, 6, "Bukti Kegiatan Harian")
	pdf.Cell(120, 6, ("Nama Pegawai: " + staf.NamaLengkap))
	pdf.Ln(-1)
	pdf.Cell(160, 6, "Pegawai RSUP Sanglah Denpasar")
	pdf.Cell(120, 6, ("NIP/Gol: " + staf.NIP + "/" + staf.GolonganPNS))
	pdf.Ln(-1)
	pdf.Cell(160, 6, ("Bulan: " + fmt.Sprintf("%s %d", months[tgl.Month()-1], tgl.Year())))
	pdf.Cell(120, 6, "Tempat Tugas: IGD RSUP Sanglah")
	pdf.Ln(-1)
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(10, 20, "No", "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 20, "Uraian", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 20, "Poin", "1", 0, "C", false, 0, "")
	pdf.CellFormat(176, 10, "Jumlah Kegiatan Harian", "1", 2, "C", false, 0, "")
	// range list iki
	for k, v := range list {
		if k < 16 {
			pdf.CellFormat(11, 10, strconv.Itoa(k+1)+shift[v.Shift], "1", 0, "C", false, 0, "")
		}
	}
	// for i := 1; i < 17; i++ {
	// 	pdf.CellFormat(11, 10, strconv.Itoa(i), "1", 0, "C", false, 0, "")
	// }
	pdf.SetXY(266, 28)
	pdf.CellFormat(25, 20, "Jumlah Poin", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(10, 24, "1", "1", 0, "C", false, 0, "")
	/////////////////////////////
	pdf.MultiCell(50, 6, "Melakukan pelayanan medik umum (per pasien : pemeriksaan rawat jalan, IGD, visite rawat inap, tim medis diskusi", "1", "L", false)
	pdf.SetXY(70, 48)
	pdf.CellFormat(20, 24, "0,0032", "1", 0, "C", false, 0, "")
	for k, v := range list {
		if k < 16 {
			if v.IKI1 == 0 && v.IKI2 == 0 {
				pdf.CellFormat(11, 24, "", "1", 0, "C", false, 0, "")
			} else {
				pdf.CellFormat(11, 24, strconv.Itoa(v.IKI1), "1", 0, "C", false, 0, "")
			}
		}
	}
	///////////////////////////////
	pdf.CellFormat(25, 24, a, "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.CellFormat(10, 12, "2", "1", 0, "C", false, 0, "")
	pdf.MultiCell(50, 6, "Melakukan tindakan medik umum tingkat sederhana (per tindakan)", "1", "L", false)
	pdf.SetXY(70, 72)
	pdf.CellFormat(20, 12, "0,01", "1", 0, "C", false, 0, "")
	///////////////////////////////
	for k, v := range list {
		if k < 16 {
			if v.IKI1 == 0 && v.IKI2 == 0 {
				pdf.CellFormat(11, 12, "", "1", 0, "C", false, 0, "")
			} else {
				pdf.CellFormat(11, 12, strconv.Itoa(v.IKI2), "1", 0, "C", false, 0, "")
			}
		}
	}
	//////////////////////////////////
	pdf.CellFormat(25, 12, b, "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.Ln(-1)
	//////////////////////////////////
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(10, 20, "No", "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 20, "Uraian", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 20, "Poin", "1", 0, "C", false, 0, "")
	pdf.CellFormat(176, 10, "Jumlah Kegiatan Harian", "1", 2, "C", false, 0, "")
	for k, v := range list {
		if k >= 16 {
			pdf.CellFormat(11, 10, strconv.Itoa(k+1)+shift[v.Shift], "1", 0, "C", false, 0, "")
		}
	}
	// for i := 17; i < 32; i++ {
	// 	pdf.CellFormat(11, 10, strconv.Itoa(i), "1", 0, "C", false, 0, "")
	// }
	pdf.SetFont("Arial", "B", 7)
	pdf.MultiCell(11, 5, "Jumlah Poin", "1", "C", false)
	pdf.SetFont("Arial", "B", 9)
	pdf.SetXY(266, 96)
	pdf.MultiCell(25, 20, "Jumlah X Poin", "1", "C", false)
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(10, 24, "1", "1", 0, "C", false, 0, "")

	pdf.MultiCell(50, 6, "Melakukan pelayanan medik umum (per pasien : pemeriksaan rawat jalan, IGD, visite rawat inap, tim medis diskusi", "1", "L", false)
	pdf.SetXY(70, 116)
	pdf.CellFormat(20, 24, "0,0032", "1", 0, "C", false, 0, "")
	for k, v := range list {
		if k >= 16 {
			if v.IKI1 == 0 && v.IKI2 == 0 {
				pdf.CellFormat(11, 24, "", "1", 0, "C", false, 0, "")
			} else {
				pdf.CellFormat(11, 24, strconv.Itoa(v.IKI1), "1", 0, "C", false, 0, "")
			}
		}
	}

	pdf.CellFormat(11, 24, c, "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 24, e, "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.CellFormat(10, 12, "2", "1", 0, "C", false, 0, "")
	pdf.MultiCell(50, 6, "Melakukan tindakan medik umum tingkat sederhana (per tindakan)", "1", "L", false)
	pdf.SetXY(70, 140)
	pdf.CellFormat(20, 12, "0,01", "1", 0, "C", false, 0, "")
	for k, v := range list {
		if k >= 16 {
			if v.IKI1 == 0 && v.IKI2 == 0 {
				pdf.CellFormat(11, 12, "", "1", 0, "C", false, 0, "")
			} else {
				pdf.CellFormat(11, 12, strconv.Itoa(v.IKI2), "1", 0, "C", false, 0, "")
			}
		}
	}

	pdf.CellFormat(11, 12, d, "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 12, f, "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.CellFormat(256, 6, "Jumlah Point X Volume kegiatan pelayanan", "1", 0, "R", false, 0, "")
	pdf.CellFormat(25, 6, g, "1", 1, "C", false, 0, "")
	pdf.CellFormat(256, 6, "Target Point kegiatan pelayanan", "1", 0, "R", false, 0, "")
	pdf.CellFormat(25, 6, fmt.Sprintf("%.4f", staf.TargetIKI), "1", 1, "C", false, 0, "")
	pdf.Ln(-1)

	for _, v := range list {
		if len(v.KegiatanLuar) != 0 {
			for m, n := range v.KegiatanLuar {
				pdf.Cell(30, 6, (strconv.Itoa(m+1) + ". " + n.Diagnosis + " (" + v.Tanggal.Format("02-01-2006") + ") "))
				pdf.Ln(-1)
				pdf.Cell(40, 6, "")
			}
		}
	}

	// Buku Catatan Pasien

	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	wd := pdf.GetStringWidth("Buku Catatan Pribadi")
	pdf.SetX((210 - wd) / 2)
	pdf.Cell(wd, 9, "Buku Catatan Pribadi")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(20, 5, "Nama")
	pdf.Cell(105, 5, (": " + staf.NamaLengkap))
	pdf.Ln(-1)
	pdf.Cell(20, 5, "Bulan")
	pdf.Cell(105, 5, (": " + fmt.Sprintf("%s %d", months[tgl.Month()-1], tgl.Year())))
	pdf.Ln(-1)
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(9, 20, "No", "1", 0, "C", false, 0, "")
	pdf.CellFormat(18, 20, "Tanggal", "1", 0, "C", false, 0, "")
	pdf.CellFormat(17, 20, "No. CM", "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, 20, "Nama", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 20, "Diagnosis", "1", 0, "C", false, 0, "")

	pdf.MultiCell(20, 5, "Melakukan pelayanan medik umum", "1", "C", false)

	pdf.SetXY(174, 35)
	pdf.MultiCell(25, 4, "Melakukan tindakan medik umum tingkat sederhana", "1", "C", false)
	var num int
	var nourut string
	for k, v := range list {
		for m, n := range list[k].Kunjungan {
			pdf.SetFont("Arial", "", 8)
			diag := properDiagnosis(n.Diagnosis)
			if len(diag) > 20 {
				diag = diag[:21]
			}
			tang := n.TglKunjungan.Format("02-01-2006")

			if k == 0 {
				nourut = strconv.Itoa(m + 1)
			} else {
				nourut = strconv.Itoa(num + m + 1)
			}
			nocm := n.NoCM
			nam := properCapital(n.NamaPasien)
			if len(nam) > 25 {
				nam = nam[:26]
			}
			pdf.CellFormat(9, 7, nourut, "1", 0, "C", false, 0, "")
			pdf.CellFormat(18, 7, tang, "1", 0, "C", false, 0, "")
			pdf.CellFormat(17, 7, nocm, "1", 0, "C", false, 0, "")
			pdf.CellFormat(60, 7, nam, "1", 0, "L", false, 0, "")
			pdf.CellFormat(40, 7, diag, "1", 0, "L", false, 0, "")
			pdf.SetFont("ZapfDingbats", "", 8)
			if n.IKI == "1" {
				pdf.CellFormat(20, 7, "4", "1", 0, "C", false, 0, "")
				pdf.CellFormat(25, 7, "", "1", 0, "C", false, 0, "")
				pdf.Ln(-1)
			} else {
				pdf.CellFormat(20, 7, "", "1", 0, "C", false, 0, "")
				pdf.CellFormat(25, 7, "4", "1", 0, "C", false, 0, "")
				pdf.Ln(-1)
			}
		}

		num = num + len(v.Kunjungan)
	}
	pdf.SetTitle(fmt.Sprintf("%s %d", months[tgl.Month()-1], tgl.Year()), true)
	t := new(bytes.Buffer)
	err := pdf.Output(t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func getResumeIKI(l []bcp) (string, string, string, string, string, string, string) {
	var a, b, c, d int

	for k, v := range l {
		switch {
		case k < 16:
			a = a + v.IKI1
			b = b + v.IKI2
		case k >= 16:
			c = c + v.IKI1
			d = d + v.IKI2
		}
	}
	f := float32(a+c) * 0.0032
	g := float32(b+d) * 0.01
	e := f + g
	return strconv.Itoa(a), strconv.Itoa(b), strconv.Itoa(a + c), strconv.Itoa(b + d), fmt.Sprintf("%.4f", f), fmt.Sprintf("%.4f", g), fmt.Sprintf("%.4f", e)
}

func properCapital(input string) string {
	words := strings.Fields(input)
	smallwords, _ := readBucket("script-igdsanglah/scripts", "nama.txt")

	for index, word := range words {
		// strings.ToLower(word)
		if strings.Contains(smallwords, " "+word+" ") {
			words[index] = word
		} else {
			words[index] = strings.Title(strings.ToLower(word))
		}
	}
	return strings.Join(words, " ")
}

func properDiagnosis(input string) string {
	words := strings.Fields(input)
	smallwords, _ := readBucket("script-igdsanglah/scripts", "diagnosis.txt")

	for index, word := range words {
		if strings.Contains(smallwords, " "+word+" ") {
			words[index] = word
		} else {
			words[index] = strings.Title(strings.ToLower(word))
		}
	}
	return strings.Join(words, " ")
}

func getDetailPasien(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/get-detail-pasien" {
		http.NotFound(w, r)
		return
	}
	js := decodeReq(r)
	r.Body.Close()

	pts := dataPasien{
		NomorCM: js.Data3,
	}

	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "igdsanglah")
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal membuat klien datastore (CM): %v", err)
		return
	}
	k := datastore.NameKey("DataPasien", pts.NomorCM, datastore.NameKey("IGD", "fasttrack", nil))
	err = client.Get(ctx, k, &pts)
	if err != nil && err != datastore.ErrNoSuchEntity {
		errorCheck(w, http.StatusInternalServerError, "Gagal mengambil data pasien: kesalahan pada server %v", err)
		return
	}
	pts.LinkID = k.Encode()
	det := detailPasien{
		IDPasien: pts,
	}

	q := datastore.NewQuery("KunjunganPasien").Ancestor(k)
	it := client.Run(ctx, q)
	var list []kunjunganPasien
	for {
		var p kunjunganPasien
		key, err := it.Next(&p)
		if err == iterator.Done {
			break
		}
		if err != nil {
			errorCheck(w, http.StatusInternalServerError, "Gagal mengambil list pasien: %v", err)
			return
		}
		p.LinkID = key.Encode()
		list = append(list, p)
	}
	det.ListKunjungan = list
	script, err := genTemplate(det, "detail-pasien")
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal membuat tabel: %v", err)
		return
	}
	w.WriteHeader(200)
	fmt.Fprint(w, script)
}

func ubahDataDokter(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ubah-data-dokter" {
		http.NotFound(w, r)
		return
	}

	js := decodeReq(r)
	r.Body.Close()

	tgl, err := time.Parse("2006-01-02", js.Data5)
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal membuat tanggal lahir: %v", err)
		return
	}
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "igdsanglah")
	if err != nil {
		log.Printf("Could not create datastore client: %v", err)
	}
	tar, err := strconv.ParseFloat(js.Data9, 32)
	if err != nil {
		errorCheck(w, http.StatusInternalServerError, "Gagal mengubah angka: %v", err)
		return
	}
	var det DetailStaf
	det.NamaLengkap = js.Data2
	det.Bagian = js.Data6
	det.TanggalLahir = tgl.In(zonaIndo())
	det.TargetIKI = tar

	if js.Data8 == "" && js.Data4 == "1" {
		det.NIP = js.Data3
		det.GolonganPNS = js.Data7

		q := datastore.NewQuery("Staff")
		it := client.Run(ctx, q)
		var ke *datastore.Key
		for {
			var st Staff
			k, err := it.Next(&st)
			if err != nil && err == iterator.Done {
				break
			}
			if st.Email == js.Data1 {
				ke = k
			}
		}
		ku, err := client.Put(ctx, datastore.NameKey("DetailStaf", "", ke), &det)
		if err != nil {
			errorCheck(w, http.StatusInternalServerError, "Gagal menyimpan perubahan dokter: %v", err)
			return
		}
		det.LinkID = ku.Encode()
	} else if js.Data8 == "" && js.Data4 != "1" {

		det.NPP = js.Data3

		q := datastore.NewQuery("Staff")
		it := client.Run(ctx, q)
		var ke *datastore.Key
		for {
			var st Staff
			k, err := it.Next(&st)
			if err != nil && err == iterator.Done {
				break
			}
			if st.Email == js.Data1 {
				ke = k
			}
		}
		ku, err := client.Put(ctx, datastore.NameKey("DetailStaf", "", ke), &det)
		if err != nil {
			errorCheck(w, http.StatusInternalServerError, "Gagal menyimpan perubahan dokter: %v", err)
			return
		}
		det.LinkID = ku.Encode()
	} else if js.Data8 != "" && js.Data4 == "1" {
		det.NIP = js.Data3
		det.GolonganPNS = js.Data7

		k, err := datastore.DecodeKey(js.Data8)
		if err != nil {
			errorCheck(w, http.StatusInternalServerError, "Gagal mendecode key: %v", err)
			return
		}
		_, err = client.Put(ctx, k, &det)
		if err != nil {
			errorCheck(w, http.StatusInternalServerError, "Gagal menyimpan perubahan dokter: %v", err)
			return
		}
		det.LinkID = k.Encode()
	} else {
		det.NPP = js.Data3

		k, err := datastore.DecodeKey(js.Data8)
		if err != nil {
			errorCheck(w, http.StatusInternalServerError, "Gagal mendecode key: %v", err)
			return
		}
		_, err = client.Put(ctx, k, &det)
		if err != nil {
			errorCheck(w, http.StatusInternalServerError, "Gagal menyimpan perubahan dokter: %v", err)
			return
		}
		det.LinkID = k.Encode()
	}

	w.WriteHeader(200)
	fmt.Fprint(w, det)
}

func getCalBali() (string, error) {
	resp, err := http.Get("http://kalenderbali.org/kbdwidget.php?id=7")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	date := string(body)
	return date[26 : len(date)-11], err
}
