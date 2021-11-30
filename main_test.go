package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

var a = App{}

func TestMain(m *testing.M) {
	a.Initialize("root", "", "localhost", "midwestemo")
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func clearTable() {
	_, err := a.DB.Exec("TRUNCATE midwestemo.news")
	if err != nil {
		log.Println(err)
	}

}

func addNews(count int) {
	if count < 1 {
		count = 1
	}

	stmtDel, err := a.DB.Prepare("INSERT INTO news(date, headline, announce, text, pic) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
	}

	defer func(stmtDel *sql.Stmt) {
		err := stmtDel.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(stmtDel)

	for i := 0; i < count; i++ {
		_, err := stmtDel.Exec(time.Now(), "Headline "+strconv.Itoa(i), "Announce "+strconv.Itoa(i), "Text "+strconv.Itoa(i), "Pic "+strconv.Itoa(i))
		if err != nil {
			log.Println(err)
		}
	}
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/news/", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	response := rr

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "null", response.Body.String())
}

func TestGetNonExistentNews(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/news/11/", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	response := rr

	assert.Equal(t, http.StatusNotFound, response.Code)

	var m map[string]string
	err := json.Unmarshal(response.Body.Bytes(), &m)
	if err != nil {
		log.Println(err)
	}
	assert.Equal(t, "News not found", m["error"])
}

func TestAddNews(t *testing.T) {
	clearTable()

	var jsonStr = []byte(`{"date":"2021-11-19 20:54:21", "headline":"test", "announce":"test", "text":"test", "pic":"test"}`)
	req, _ := http.NewRequest("POST", "/news/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	response := rr

	assert.Equal(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	err := json.Unmarshal(response.Body.Bytes(), &m)
	if err != nil {
		log.Println(err)
	}

	assert.Equal(t, 0.0, m["id"])
	assert.Equal(t, "2021-11-19 20:54:21", m["date"])
	assert.Equal(t, "test", m["headline"])
	assert.Equal(t, "test", m["announce"])
	assert.Equal(t, "test", m["text"])
	assert.Equal(t, "test", m["pic"])
}

func TestGetNewsById(t *testing.T) {
	clearTable()
	addNews(1)

	req, _ := http.NewRequest("GET", "/news/1/", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	response := rr

	assert.Equal(t, http.StatusOK, response.Code)
}

func TestDeleteNewsById(t *testing.T) {
	clearTable()
	addNews(1)

	req, _ := http.NewRequest("GET", "/news/1/", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	response := rr

	assert.Equal(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/news/1/", nil)
	rr = httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	response = rr

	assert.Equal(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/news/1/", nil)
	rr = httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	response = rr

	assert.Equal(t, http.StatusNotFound, response.Code)
}

func TestGetAllReleases(t *testing.T) {
	clearTable()
	addNews(1)

	req, _ := http.NewRequest("GET", "/releases/", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	response := rr

	assert.Equal(t, http.StatusOK, response.Code)
}
