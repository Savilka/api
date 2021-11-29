package main

import (
	"api/internal/store"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"strconv"
)

type apiServer struct {
	store *store.Store
}

func NewApiServer() *apiServer {
	apiStore := store.New()
	return &apiServer{apiStore}
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		return
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	server := NewApiServer()

	router.HandleFunc("/news/", server.getAllNewsHandler).Methods("GET")
	router.HandleFunc("/news/{id}", server.getNewsByIdHandler).Methods("GET")
	router.HandleFunc("/releases/", server.getAllReleasesHandler).Methods("GET")
	router.HandleFunc("/news/", server.addNewsHandler).Methods("POST")
	router.HandleFunc("/news/{id}", server.deleteNewsById).Methods("DELETE")

	handler := cors.Default().Handler(router)

	//Запуск сервера
	log.Fatal(http.ListenAndServe("localhost:10000", handler))

}

func (as *apiServer) getAllNewsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling get all news at %s\n", r.URL.Path)

	allNews, err := as.store.GetAllNews()
	if err != nil {
		log.Println(err.Error())
	}

	renderJSON(w, allNews)

}

func (as *apiServer) getAllReleasesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling get all news at %s\n", r.URL.Path)

	allReleases, err := as.store.GetAllReleases()
	if err != nil {
		log.Println(err.Error())
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Origin", "Content-Type")
	renderJSON(w, allReleases)

}

func (as *apiServer) getNewsByIdHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	log.Printf("handling get news with id = '%d' at %s\n", id, r.URL.Path)

	task, err := as.store.GetNewsById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, task)
}

func (as *apiServer) deleteNewsById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	log.Printf("handling delete news with id = '%d' at %s\n", id, r.URL.Path)

	err := as.store.DeleteNewsById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func (as *apiServer) addNewsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling add news at %s\n", r.URL.Path)

	if r.Body == nil {

	}
}
