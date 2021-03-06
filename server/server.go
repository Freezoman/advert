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
	CreateAdvert(data *multipart.Form) (string, error)
}

func (s *Server) getAdvert(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/advert/"):]
	keys, ok := r.URL.Query()["fields"]
	fields := []string{}
	if ok {
		fields = append(fields, strings.Split(keys[0], ",")...)
	}
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

	res, err := s.app.CreateAdvert(formData)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(res))
	return
}

func (s *Server) ServerStart() {
	http.HandleFunc("/advert/", s.getAdvert)
	http.HandleFunc("/adverts", s.getAdverts)
	http.HandleFunc("/upload", s.uploadAdvert)
	http.Handle("/files/", http.StripPrefix("/files", http.FileServer(http.Dir("files"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
