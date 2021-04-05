package live

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	default:
		return fmt.Errorf("Video request %s not supported", reqReq.Name)
	}
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
