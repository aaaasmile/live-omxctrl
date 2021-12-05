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

func handleVideoRequest(w http.ResponseWriter, req *http.Request) error {
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
	case "FetchVideo":
		return fetchVideoReq(rawbody, w, req)
	case "ScanVideo":
		return scanVideoReq(rawbody, w, req)
	default:
		return fmt.Errorf("Video request %s not supported", reqReq.Name)
	}
}

func scanVideoReq(rawbody []byte, w http.ResponseWriter, req *http.Request) error {
	start := time.Now()
	videoPath := conf.Current.VideoDir
	list, err := getVideoFiles(videoPath)
	if err != nil {
		return err
	}
	log.Println("Video file found: ", len(list))

	trxdelete, err := liteDB.GetNewTransaction()
	if err != nil {
		return err
	}

	err = liteDB.DeleteAllVideo(trxdelete)
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

	err = liteDB.InsertVideoList(trx, list)
	if err != nil {
		return err
	}
	err = trx.Commit()
	if err != nil {
		return err
	}

	log.Println("Scan and store processing time ", time.Now().Sub(start))

	return fetchVideoReq(rawbody, w, req)
}

func fetchVideoReq(rawbody []byte, w http.ResponseWriter, req *http.Request) error {
	paraReq := struct {
		PageIx   int `json:"pageix"`
		PageSize int `json:"pagesize"`
	}{}
	if err := json.Unmarshal(rawbody, &paraReq); err != nil {
		return err
	}
	log.Println("video Request", paraReq)
	list, err := liteDB.FetchVideo(paraReq.PageIx, paraReq.PageSize)
	if err != nil {
		return err
	}

	type ItemRes struct {
		ID          int    `json:"id"`
		Type        string `json:"type"`
		Title       string `json:"title"`
		URI         string `json:"uri"`
		DurationStr string `json:"durationstr"`
	}

	res := struct {
		Video  []ItemRes `json:"video"`
		PageIx int       `json:"pageix"`
	}{
		Video:  make([]ItemRes, 0),
		PageIx: paraReq.PageIx,
	}
	for _, item := range list {
		pp := ItemRes{
			ID:          item.ID,
			Type:        item.Type,
			Title:       item.Title,
			URI:         item.URI,
			DurationStr: item.Duration,
		}
		res.Video = append(res.Video, pp)
	}

	return writeResponseNoWsBroadcast(w, res)
}

func getVideoFiles(rootPath string) ([]*db.ResUriItem, error) {
	rootPath, _ = filepath.Abs(rootPath)
	arr := []*db.ResUriItem{}
	filterVideo := []string{".mp4", ".avi", ".mkv"}
	log.Printf("Process path %s", rootPath)
	if info, err := os.Stat(rootPath); err == nil && info.IsDir() {
		arr, err = getVideosinDir(rootPath, filterVideo, arr)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return arr, nil
}

func getVideosinDir(dirAbs string, filterVideo []string, parentItems []*db.ResUriItem) ([]*db.ResUriItem, error) {
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
			r, err = getVideosinDir(pathAbsItem, filterVideo, r)
			if err != nil {
				return nil, err
			}
		} else {
			//fmt.Println("** file is ", f.Name())
			ext := filepath.Ext(pathAbsItem)
			for _, v := range filterVideo {
				if v == ext {
					item := db.ResUriItem{
						URI:   pathAbsItem,
						Title: strings.ReplaceAll(f.Name(), v, ""),
						Type:  strings.Trim(ext, "."),
					}

					r = append(r, &item)
					break
				}
			}
		}
	}
	return r, nil
}
