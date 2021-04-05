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

	"github.com/aaaasmile/live-omxctrl/conf"
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
	videoPath := conf.Current.VideoDir
	list, err := getVideoFiles(videoPath)
	if err != nil {
		return err
	}
	log.Println("Video file found: ", len(list), list)

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

func getVideoFiles(rootPath string) ([]string, error) {
	rootPath, _ = filepath.Abs(rootPath)
	onlyFiles := []string{}
	filterVideo := []string{".mp4", ".avi"}
	log.Printf("Process path %s", rootPath)
	if info, err := os.Stat(rootPath); err == nil && info.IsDir() {
		arr := []string{}
		arr, err = getFilesinDir(rootPath, filterVideo, arr)
		if err != nil {
			return nil, err
		}
		//fmt.Println("Dir process result: ", arr)
		for _, ele := range arr {
			onlyFiles = append(onlyFiles, ele)
		}
	} else {
		return nil, err
	}

	return onlyFiles, nil
}

func getFilesinDir(dirAbs string, filterVideo []string, parentFiles []string) ([]string, error) {
	r := parentFiles
	log.Println("Scan dir ", dirAbs)
	files, err := ioutil.ReadDir(dirAbs)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		itemAbs := path.Join(dirAbs, f.Name())
		if info, err := os.Stat(itemAbs); err == nil && info.IsDir() {
			//fmt.Println("** Sub dir found ", f.Name())
			r, err = getFilesinDir(itemAbs, filterVideo, r)
			if err != nil {
				return nil, err
			}
		} else {
			//fmt.Println("** file is ", f.Name())
			ext := filepath.Ext(itemAbs)
			for _, v := range filterVideo {
				if v == ext {
					r = append(r, itemAbs)
					break
				}
			}
		}
	}
	return r, nil
}
