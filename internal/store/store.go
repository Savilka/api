package store

import (
	"database/sql"
	"log"
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

func (n *News) GetNewsById(db *sql.DB) error {
	stmtOut, err := db.Prepare("SELECT * FROM `news` WHERE id = ?")
	if err != nil {
		return err
	}

	defer func(stmtOut *sql.Stmt) {
		err := stmtOut.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(stmtOut)

	err = stmtOut.QueryRow(n.ID).Scan(&n.ID, &n.Date, &n.Headline, &n.Announce, &n.Text, &n.Pic)
	if err != nil {
		return err
	}

	return nil
}

func GetAllNews(db *sql.DB) ([]News, error) {
	results, err := db.Query("SELECT * FROM `news`")
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

func GetAllReleases(db *sql.DB) ([]Release, error) {
	results, err := db.Query("SELECT * FROM `releases`")
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

func (n *News) DeleteNewsById(db *sql.DB) error {
	stmtDel, err := db.Prepare("DELETE FROM `news` WHERE id = ?")
	if err != nil {
		return err
	}

	defer func(stmtDel *sql.Stmt) {
		err := stmtDel.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(stmtDel)

	_, err = stmtDel.Exec(n.ID)

	return err
}

func (n *News) AddNews(db *sql.DB) error {
	stmtDel, err := db.Prepare("INSERT INTO news(date, headline, announce, text, pic) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer func(stmtDel *sql.Stmt) {
		err := stmtDel.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(stmtDel)

	res, err := stmtDel.Exec(n.Date, n.Headline, n.Announce, n.Text, n.Pic)
	if err != nil {
		return err
	}

	lid, err := res.LastInsertId()
	if err != nil {
		return err
	}

	n.ID = int(lid)

	return nil
}
