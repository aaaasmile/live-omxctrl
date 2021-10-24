package live

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func handleRadioRequest(w http.ResponseWriter, req *http.Request) error {
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
	case "FetchRadio":
		return fetchRadioReq(rawbody, w, req)
	default:
		return fmt.Errorf("Radio request %s not supported", reqReq.Name)
	}
}

func fetchRadioReq(rawbody []byte, w http.ResponseWriter, req *http.Request) error {
	paraReq := struct {
		PageIx   int `json:"pageix"`
		PageSize int `json:"pagesize"`
	}{}
	if err := json.Unmarshal(rawbody, &paraReq); err != nil {
		return err
	}
	log.Println("radio Request", paraReq)
	list, err := liteDB.FetchRadio(paraReq.PageIx, paraReq.PageSize)
	if err != nil {
		return err
	}

	type ItemRes struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Genre       string `json:"genre"`
		URI         string `json:"uri"`
	}

	res := struct {
		Radio  []ItemRes `json:"radio"`
		PageIx int       `json:"pageix"`
	}{
		Radio:  make([]ItemRes, 0),
		PageIx: paraReq.PageIx,
	}
	for _, item := range list {
		pp := ItemRes{
			ID:          item.ID,
			Title:       item.Title,
			URI:         item.URI,
			Description: item.Description,
			Genre:       item.Genre,
		}
		res.Radio = append(res.Radio, pp)
	}

	return writeResponseNoWsBroadcast(w, res)
}
