package main

import (
	"encoding/json"
	"errors"
	"ex04/db"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Store interface {
	GetPlaces(limit int, offset int) ([]db.Place, int, error)
	GetPlacesRecommed(lat, lon float64) ([]db.Place, error)
}

var (
	base      Store
	err       error
	secretKey = []byte("gr57hgfLHJ++rt*EE")
)

func init() {
	base, err = db.NewElasticsearch()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/api/places", apiPlaces)
	http.HandleFunc("/api/recommend", apiRecommend)
	http.HandleFunc("/api/get_token", getToken)
	log.Fatal(http.ListenAndServe(":8888", nil))
}

type Data struct {
	Total  int        `json:"total"`
	Places []db.Place `json:"places"`
	Page   int        `json:"page"`
	Last   int        `json:"last"`
}

func home(w http.ResponseWriter, r *http.Request) {
	var res Data
	pageStr := r.URL.Query().Get("page")
	if res.Page, err = strconv.Atoi(pageStr); err != nil {
		http.Error(w, fmt.Sprintf("error : Invalid 'page' value: '%s'", pageStr), http.StatusBadRequest)
		return
	}

	limit := 10
	offset := (res.Page - 1) * limit

	if res.Places, res.Total, err = base.GetPlaces(limit, offset); err != nil {
		log.Println(err)
		return
	}
	res.Last = int(math.Ceil(float64(res.Total) / float64(limit)))
	tmpl, err := template.New("index.html").Funcs(
		template.FuncMap{
			"sum": sum,
			"sub": sub,
		},
	).ParseFiles("template/index.html")

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, res)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func apiPlaces(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res Data
	pageStr := r.URL.Query().Get("page")
	if res.Page, err = strconv.Atoi(pageStr); err != nil {
		http.Error(w, fmt.Sprintf("\"error\" : Invalid 'page' value: '%s'", pageStr), http.StatusBadRequest)
		return
	}
	limit := 10
	offset := (res.Page - 1) * limit

	if res.Places, res.Total, err = base.GetPlaces(limit, offset); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res.Last = int(math.Ceil(float64(res.Total) / float64(limit)))
	response, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func sum(x, y int) int {
	return x + y
}

func sub(x, y int) int {
	return x - y
}

func apiRecommend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	str := r.Header.Get("Authorization")
	token := strings.Split(str, " ")
	if len(token) != 2 {
		http.Error(w, "empty jwt token", http.StatusUnauthorized)
		return
	}
	if err = parseToken(token[1]); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	if err != nil {
		http.Error(w, "Invalid 'lat' number", http.StatusBadRequest)
		return
	}
	lon, err := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
	if err != nil {
		http.Error(w, "Invalid 'lon' number", http.StatusBadRequest)
		return
	}
	pl, err := base.GetPlacesRecommed(lat, lon)
	if err != nil {
		log.Println(err)
		return
	}
	res := struct {
		Name   string     `json:"name"`
		Places []db.Place `json:"places"`
	}{
		Name:   "Recommendation",
		Places: pl,
	}
	response, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func getToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token, err := generateJWT()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	response, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func generateJWT() (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(12 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func parseToken(accessToken string) error {
	_, err := jwt.Parse(
		accessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(secretKey), nil
		},
	)
	if err != nil {
		return err
	}
	return nil
}
