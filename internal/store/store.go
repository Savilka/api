package store

import (
	"database/sql"
	"log"
	"time"
)

// News -  структура, описывающая одну 'новость'
type News struct {
	ID       int    `json:"id"`
	Date     string `json:"date"`
	Headline string `json:"headline"`
	Announce string `json:"announce"`
	Text     string `json:"text"`
	Pic      string `json:"pic"`
}

// Release -  структура, описывающая однин 'релиз'
type Release struct {
	ID     int    `json:"id"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Pic    string `json:"pic"`
}

type Store struct {
	db *sql.DB
}

func New() *Store {
	s := &Store{}
	var err error
	s.db, err = sql.Open("mysql", "yfnk95pabxc7o8pb:cji3ywltts8o22ui"+
		"@tcp(n2o93bb1bwmn0zle.chr7pe7iynqr.eu-west-1.rds.amazonaws.com:3306)/qytxwbkkn14vj0yd")
	if err != nil {
		panic(err.Error())
	}

	return s
}

func (s *Store) GetNewsById(id int) (News, error) {
	stmtOut, err := s.db.Prepare("SELECT * FROM `news` WHERE id = ?")
	if err != nil {
		return News{}, err
	}

	defer func(stmtOut *sql.Stmt) {
		err := stmtOut.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(stmtOut)

	var news News
	err = stmtOut.QueryRow(id).Scan(&news.ID, &news.Date, &news.Headline, &news.Announce, &news.Text, &news.Pic)
	if err != nil {
		return News{}, err
	}

	return news, nil
}

func (s *Store) GetAllNews() ([]News, error) {
	results, err := s.db.Query("SELECT * FROM `news`")
	if err != nil {
		return []News{}, err
	}

	var allNews []News

	for results.Next() {
		var nextNews News
		err = results.Scan(&nextNews.ID, &nextNews.Date, &nextNews.Headline, &nextNews.Announce, &nextNews.Text, &nextNews.Pic)
		if err != nil {
			return []News{}, err
		}
		allNews = append(allNews, nextNews)
	}

	return allNews, nil
}

func (s *Store) GetAllReleases() ([]Release, error) {
	results, err := s.db.Query("SELECT * FROM `releases`")
	if err != nil {
		return []Release{}, err
	}

	var allReleases []Release

	for results.Next() {
		var nextRelease Release
		err = results.Scan(&nextRelease.ID, &nextRelease.Artist, &nextRelease.Album, &nextRelease.Pic)
		if err != nil {
			return []Release{}, err
		}
		allReleases = append(allReleases, nextRelease)
	}

	return allReleases, nil
}

func (s *Store) DeleteNewsById(id int) (err error) {
	stmtDel, err := s.db.Prepare("DELETE FROM `news` WHERE id = ?")
	if err != nil {
		return err
	}

	defer func(stmtDel *sql.Stmt) {
		err := stmtDel.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(stmtDel)

	_, err = stmtDel.Query(id)

	return err
}

func (s *Store) AddNews(date time.Time, headline string, announce string, text string, pic string) (err error) {
	stmtDel, err := s.db.Prepare("INSERT INTO news(date, headline, announce, text, pic) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer func(stmtDel *sql.Stmt) {
		err := stmtDel.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(stmtDel)

	_, err = stmtDel.Query(date, headline, announce, text, pic)

	return err
}
