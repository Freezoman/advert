package server

import (
	"fmt"
	"log"
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
	GetJsonAdvert(id int, optFields ...string) string
	GetJsonAdverts(sortField string, order string) string
}

func (s *Server) getAdvert(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["fields"]
	fields := []string{}
	if ok {
		fields = append(fields, strings.Split(keys[0], ",")...)
	}
	// fmt.Println(fields)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, s.app.GetJsonAdvert(1, fields...))
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
func adv(w http.ResponseWriter, r *http.Request) {

}
func advs(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) ServerStart() {
	http.HandleFunc("/advert", s.getAdvert)
	http.HandleFunc("/adverts", s.getAdverts)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
