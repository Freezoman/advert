package server

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

type Server struct {
	app IApp
}

func Init(a IApp) *Server {
	return &Server{
		app: a,
	}
}

type IApp interface {
	GetJsonAdvert(id string, optFields ...string) string
	GetJsonAdverts(sortField string, order string) string
	CreateAdvert(data *multipart.Form) error
}

func (s *Server) getAdvert(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/advert/"):]
	keys, ok := r.URL.Query()["fields"]
	fields := []string{}
	if ok {
		fields = append(fields, strings.Split(keys[0], ",")...)
	}
	// fmt.Println(fields)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, s.app.GetJsonAdvert(id, fields...))
}
func (s *Server) getAdverts(w http.ResponseWriter, r *http.Request) {
	field, order := "", ""
	keys, ok := r.URL.Query()["field"]
	if ok {
		field = keys[0]
	}
	keys, ok = r.URL.Query()["order"]
	if ok {
		order = keys[0]
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(s.app.GetJsonAdverts(field, order)))
}
func (s *Server) uploadAdvert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write([]byte("404 could not find page"))
		return
	}
	r.ParseMultipartForm(10 << 20)
	formData := r.MultipartForm

	err := s.app.CreateAdvert(formData)
	if err == nil {
		w.Write([]byte("advert can't be created"))
		return
	}
	w.Write([]byte("advert has been created successfully"))
	return
}

// func (s *Server) adv(w http.ResponseWriter, r *http.Request) {
// 	title := r.URL.Path[len("/advert/") :]
// 	fmt.Fprintf(w, title)
// }

func (s *Server) ServerStart() {
	http.HandleFunc("/advert/", s.getAdvert)
	http.HandleFunc("/adverts", s.getAdverts)
	http.HandleFunc("/upload", s.uploadAdvert)

	// http.HandleFunc("/advert/", s.adv)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
