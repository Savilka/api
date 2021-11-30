package main

import (
	"api/internal/store"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"io"
	"log"
	"net/http"
	"strconv"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, host, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, password, host, dbname)

	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter().StrictSlash(true)

	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, cors.Default().Handler(a.Router)))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/news/", a.getAllNewsHandler).Methods("GET")
	a.Router.HandleFunc("/news/{id}/", a.getNewsByIdHandler).Methods("GET")
	a.Router.HandleFunc("/releases/", a.getAllReleasesHandler).Methods("GET")
	a.Router.HandleFunc("/news/", a.postNewsHandler).Methods("POST")
	a.Router.HandleFunc("/news/{id}/", a.deleteNewsById).Methods("DELETE")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

func (a *App) getAllNewsHandler(w http.ResponseWriter, _ *http.Request) {
	news, err := store.GetAllNews(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, news)
}

func (a *App) getAllReleasesHandler(w http.ResponseWriter, _ *http.Request) {
	releases, err := store.GetAllReleases(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, releases)
}

func (a *App) getNewsByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	n := store.News{ID: id}
	if err := n.GetNewsById(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "News not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, n)
}

func (a *App) deleteNewsById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	n := store.News{ID: id}
	if err := n.DeleteNewsById(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) postNewsHandler(w http.ResponseWriter, r *http.Request) {
	var n store.News
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&n); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(r.Body)

	if err := n.DeleteNewsById(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, n)
}
