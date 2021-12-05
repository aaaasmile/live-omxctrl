package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aaaasmile/live-omxctrl/util"
	_ "github.com/mattn/go-sqlite3"
)

type LiteDB struct {
	connDb       *sql.DB
	DebugSQL     bool
	SqliteDBPath string
}

type ResUriItem struct {
	ID                                int
	URI, Title, Description, Duration string
	Timestamp                         time.Time
	PlayPosition                      int
	DurationInSec                     int
	Type                              string
	Genre                             string
}

type ResMusicItem struct {
	Timestamp     time.Time
	ID            int
	FileOrFolder  int
	Title         string
	URI           string
	Description   string
	DurationInSec int
	ParentFolder  string
	MetaAlbum     string
	MetaArtist    string
}

func (ld *LiteDB) OpenSqliteDatabase() error {
	var err error
	dbname := util.GetFullPath(ld.SqliteDBPath)
	if _, err := os.Stat(dbname); err != nil {
		return err
	}
	log.Println("Using the sqlite file: ", dbname)
	ld.connDb, err = sql.Open("sqlite3", dbname)
	if err != nil {
		return err
	}
	return nil
}

func (ld *LiteDB) FetchVideo(pageIx int, pageSize int) ([]ResUriItem, error) {
	q := `SELECT id,Timestamp,URI,Title,Description,Duration,PlayPosition,DurationInSec,Type
		  FROM Video
		  ORDER BY Title DESC 
		  LIMIT %d OFFSET %d;`
	offsetRows := pageIx * pageSize
	q = fmt.Sprintf(q, pageSize, offsetRows)
	if ld.DebugSQL {
		log.Println("Query is", q)
	}

	rows, err := ld.connDb.Query(q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	res := make([]ResUriItem, 0)
	var tss int64
	for rows.Next() {
		item := ResUriItem{}
		tss = 0
		if err := rows.Scan(&item.ID, &tss, &item.URI, &item.Title,
			&item.Description, &item.Duration, &item.PlayPosition,
			&item.DurationInSec, &item.Type); err != nil {
			return nil, err
		}
		item.Timestamp = time.Unix(tss, 0)
		res = append(res, item)
	}
	return res, nil
}

func (ld *LiteDB) FetchMusic(parent string) ([]ResMusicItem, error) {
	q := `SELECT id,Timestamp,URI,Title,Description,DurationInSec,ParentFolder,FileOrFolder,MetaAlbum,MetaArtist
		  FROM MusicFile
		  WHERE ParentFolder=?
		  ORDER BY Title ASC;`
	if ld.DebugSQL {
		log.Println("Query is", q)
	}

	rows, err := ld.connDb.Query(q, parent)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	res := make([]ResMusicItem, 0)
	var tss int64
	for rows.Next() {
		item := ResMusicItem{}
		tss = 0
		if err := rows.Scan(&item.ID, &tss, &item.URI, &item.Title,
			&item.Description, &item.DurationInSec,
			&item.ParentFolder, &item.FileOrFolder, &item.MetaAlbum, &item.MetaArtist); err != nil {
			return nil, err
		}
		item.Timestamp = time.Unix(tss, 0)
		res = append(res, item)
	}
	return res, nil
}

func (ld *LiteDB) FetchRadioFromURI(uri string) (*ResUriItem, error) {
	q := `SELECT id,URI,Name,Description,Genre
		  FROM Radio
		  WHERE URI = "%s"
		  LIMIT 1;`
	q = fmt.Sprintf(q, uri)
	if ld.DebugSQL {
		log.Println("Query is", q)
	}
	res := ResUriItem{}

	rows, err := ld.connDb.Query(q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&res.ID, &res.URI, &res.Title, &res.Description, &res.Genre); err != nil {
			return nil, err
		}
		break
	}

	return &res, nil
}

func (ld *LiteDB) FetchRadio(pageIx int, pageSize int) ([]ResUriItem, error) {
	q := `SELECT id,URI,Name,Description,Genre
		  FROM Radio
		  ORDER BY Name DESC 
		  LIMIT %d OFFSET %d;`
	offsetRows := pageIx * pageSize
	q = fmt.Sprintf(q, pageSize, offsetRows)
	if ld.DebugSQL {
		log.Println("Query is", q)
	}

	rows, err := ld.connDb.Query(q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	res := make([]ResUriItem, 0)
	for rows.Next() {
		item := ResUriItem{}
		if err := rows.Scan(&item.ID, &item.URI, &item.Title,
			&item.Description, &item.Genre); err != nil {
			return nil, err
		}
		res = append(res, item)
	}
	return res, nil
}

func (ld *LiteDB) GetNewTransaction() (*sql.Tx, error) {
	tx, err := ld.connDb.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (ld *LiteDB) InsertVideoList(tx *sql.Tx, list []*ResUriItem) error {
	for _, item := range list {
		q := `INSERT INTO Video(Timestamp,URI,Title,Description,Duration,PlayPosition,DurationInSec,Type) VALUES(?,?,?,?,?,?,?,?);`
		if ld.DebugSQL {
			log.Println("Query is", q)
		}

		stmt, err := ld.connDb.Prepare(q)
		if err != nil {
			return err
		}

		now := time.Now()
		sqlres, err := tx.Stmt(stmt).Exec(now.Local().Unix(), item.URI, item.Title, item.Description,
			item.Duration, 0, item.DurationInSec, item.Type)
		if err != nil {
			return err
		}
		log.Println("video inserted: ", item.Title, sqlres)
	}
	return nil
}

func (ld *LiteDB) InsertMusicList(tx *sql.Tx, list []*ResMusicItem) error {
	for _, item := range list {
		q := `INSERT INTO MusicFile(Timestamp,URI,Title,Description,DurationInSec,FileOrFolder,ParentFolder,MetaAlbum,MetaArtist) VALUES(?,?,?,?,?,?,?,?,?);`
		if ld.DebugSQL {
			log.Println("Query is", q)
		}

		stmt, err := ld.connDb.Prepare(q)
		if err != nil {
			return err
		}

		now := time.Now()
		sqlres, err := tx.Stmt(stmt).Exec(now.Local().Unix(), item.URI, item.Title, item.Description,
			item.DurationInSec, item.FileOrFolder, item.ParentFolder, item.MetaAlbum, item.MetaArtist)
		if err != nil {
			return err
		}
		log.Println("video inserted: ", item.Title, sqlres)
	}
	return nil
}

func (ld *LiteDB) DeleteAllVideo(tx *sql.Tx) error {
	q := fmt.Sprintf(`DELETE FROM Video;`)
	if ld.DebugSQL {
		log.Println("Query is", q)
	}

	stmt, err := ld.connDb.Prepare(q)
	if err != nil {
		return err
	}

	_, err = tx.Stmt(stmt).Exec()
	return err
}

func (ld *LiteDB) DeleteAllMusicFiles(tx *sql.Tx) error {
	q := fmt.Sprintf(`DELETE FROM MusicFile;`)
	if ld.DebugSQL {
		log.Println("Query is", q)
	}

	stmt, err := ld.connDb.Prepare(q)
	if err != nil {
		return err
	}

	_, err = tx.Stmt(stmt).Exec()
	return err
}

func (ld *LiteDB) FetchHistory(pageIx int, pageSize int) ([]ResUriItem, error) {
	q := `SELECT id,Timestamp,URI,Title,Description,Duration,PlayPosition,DurationInSec,Type
		  FROM History
		  ORDER BY Timestamp DESC 
		  LIMIT %d OFFSET %d;`
	offsetRows := pageIx * pageSize
	q = fmt.Sprintf(q, pageSize, offsetRows)
	if ld.DebugSQL {
		log.Println("Query is", q)
	}

	rows, err := ld.connDb.Query(q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	res := make([]ResUriItem, 0)
	var tss int64
	for rows.Next() {
		item := ResUriItem{}
		tss = 0
		if err := rows.Scan(&item.ID, &tss, &item.URI, &item.Title,
			&item.Description, &item.Duration, &item.PlayPosition,
			&item.DurationInSec, &item.Type); err != nil {
			return nil, err
		}
		item.Timestamp = time.Unix(tss, 0)
		res = append(res, item)
	}
	return res, nil
}

func (ld *LiteDB) CreateHistory(uri, title, description, duration string, durinsec int, tt string) error {
	item := ResUriItem{
		URI:           uri,
		Title:         title,
		Description:   description,
		Duration:      duration,
		DurationInSec: durinsec,
		Type:          tt,
	}
	return ld.InsertHistoryItem(&item)
}

func (ld *LiteDB) InsertHistoryItem(item *ResUriItem) error {
	q := `INSERT INTO History(Timestamp,URI,Title,Description,Duration,PlayPosition,DurationInSec,Type) VALUES(?,?,?,?,?,?,?,?);`
	if ld.DebugSQL {
		log.Println("Query is", q)
	}

	stmt, err := ld.connDb.Prepare(q)
	if err != nil {
		return err
	}

	now := time.Now()
	sqlres, err := stmt.Exec(now.Local().Unix(), item.URI, item.Title, item.Description,
		item.Duration, 0, item.DurationInSec, item.Type)
	if err != nil {
		return err
	}
	log.Println("History inserted: ", sqlres)
	return nil
}
