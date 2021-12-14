package live

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/aaaasmile/live-omxctrl/conf"
	"github.com/aaaasmile/live-omxctrl/db"
)

func handleMusicRequest(w http.ResponseWriter, req *http.Request) error {
	rawbody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	reqReq := struct {
		Name string `json:"name"`
	}{}
	if err := json.Unmarshal(rawbody, &reqReq); err != nil {
		return err
	}

	switch reqReq.Name {
	case "ScanMusic":
		return scanMusicReq(rawbody, w, req)
	case "FetchMusic":
		return fetchMusicReq(rawbody, w, req)
	default:
		return fmt.Errorf("Music request %s not supported", reqReq.Name)
	}
}

func scanMusicReq(rawbody []byte, w http.ResponseWriter, req *http.Request) error {
	start := time.Now()
	musicDir := conf.Current.MusicDir
	list, err := getMusicFiles(musicDir)
	if err != nil {
		return err
	}
	log.Println("Video file found: ", len(list))

	trxdelete, err := liteDB.GetNewTransaction()
	if err != nil {
		return err
	}

	err = liteDB.DeleteAllMusicFiles(trxdelete)
	if err != nil {
		return err
	}
	err = trxdelete.Commit()
	if err != nil {
		return err
	}

	trx, err := liteDB.GetNewTransaction()
	if err != nil {
		return err
	}

	err = liteDB.InsertMusicList(trx, list)
	if err != nil {
		return err
	}
	err = trx.Commit()
	if err != nil {
		return err
	}

	log.Println("Scan and store processing time ", time.Since(start))

	return fetchMusicReq(rawbody, w, req)
}

func fetchMusicReq(rawbody []byte, w http.ResponseWriter, req *http.Request) error {
	paraReq := struct {
		Parent string `json:"parent"`
	}{}
	if err := json.Unmarshal(rawbody, &paraReq); err != nil {
		return err
	}
	log.Println("Music Request", paraReq)
	list, err := liteDB.FetchMusic(paraReq.Parent)
	if err != nil {
		return err
	}

	type MusicItemRes struct {
		ID           int    `json:"id"`
		FileOrFolder int    `json:"fileorfolder"`
		Title        string `json:"title"`
		URI          string `json:"uri"`
		DurationStr  string `json:"durationstr"`
		MetaAlbum    string `json:"metaalbum"`
		MetaArtist   string `json:"metaartist"`
	}

	res := struct {
		Music  []MusicItemRes `json:"music"`
		Parent string         `json:"parent"`
	}{
		Music:  make([]MusicItemRes, 0),
		Parent: paraReq.Parent,
	}
	for _, item := range list {
		pp := MusicItemRes{
			ID:           item.ID,
			FileOrFolder: item.FileOrFolder,
			Title:        item.Title,
			URI:          item.URI,
			DurationStr:  fmt.Sprintf("%s", time.Duration(item.DurationInSec*1000000000)),
			MetaAlbum:    item.MetaAlbum,
			MetaArtist:   item.MetaArtist,
		}
		res.Music = append(res.Music, pp)
	}

	return writeResponseNoWsBroadcast(w, res)
}

func getMusicFiles(rootPath string) ([]*db.ResMusicItem, error) {
	rootPath, _ = filepath.Abs(rootPath)
	arr := []*db.ResMusicItem{}
	filterMusic := []string{".mp3", ".ogg", ".wav"}
	log.Printf("Process path %s", rootPath)
	if info, err := os.Stat(rootPath); err == nil && info.IsDir() {
		arr, err = getMusicsInDir(rootPath, rootPath, filterMusic, arr)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return arr, nil
}

func getMusicsInDir(rootPath, dirAbs string, filterMusic []string, parentItems []*db.ResMusicItem) ([]*db.ResMusicItem, error) {
	r := parentItems
	log.Println("Scan dir ", dirAbs)
	files, err := ioutil.ReadDir(dirAbs)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		pathAbsItem := path.Join(dirAbs, f.Name())
		if info, err := os.Stat(pathAbsItem); err == nil && info.IsDir() {
			//fmt.Println("** Sub dir found ", f.Name())
			item := db.ResMusicItem{
				URI:          pathAbsItem,
				Title:        f.Name(),
				FileOrFolder: 0,
				ParentFolder: strings.Replace(dirAbs, rootPath, "", 1),
			}
			r = append(r, &item)
			r, err = getMusicsInDir(rootPath, pathAbsItem, filterMusic, r)
			if err != nil {
				return nil, err
			}
		} else {
			//fmt.Println("** file is ", f.Name())
			ext := filepath.Ext(pathAbsItem)
			for _, v := range filterMusic {
				if v == ext {
					item := db.ResMusicItem{
						URI:          pathAbsItem,
						FileOrFolder: 1,
						Title:        strings.ReplaceAll(f.Name(), v, ""),
						ParentFolder: strings.Replace(dirAbs, rootPath, "", 1),
					}

					r = append(r, &item)
					break
				}
			}
		}
	}
	return r, nil
}
