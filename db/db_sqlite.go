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
