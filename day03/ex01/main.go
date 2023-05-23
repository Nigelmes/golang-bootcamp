package main

import (
	"ex01/db"
	"github.com/elastic/go-elasticsearch/v7"
	"html/template"
	"log"
	"net/http"
)

type Store interface {
	GetPlaces(limit int, offset int) ([]db.Place, int, error)
}

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	db := db.NewElasticsearch(es)
	http.HandleFunc(
		"/", func(w http.ResponseWriter, r *http.Request) {
			var (
				res Data
				err error
			)
			res.Places, res.Total, err = db.GetPlaces(10, 10)
			tmpl, err := template.ParseFiles("template/index.html")
			if err != nil {
				http.Error(w, "Internal Server Error1", http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, res)
			if err != nil {
				log.Println(err)
				http.Error(w, "Internal Server Error2", http.StatusInternalServerError)
				return
			}
		},
	)
	http.ListenAndServe(":8888", nil)

}

type Data struct {
	Places []db.Place
	Total  int
}
