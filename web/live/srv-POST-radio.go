package live

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aaaasmile/live-omxctrl/db"
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
	case "InsertRadio":
		return insertRadioReq(rawbody, w, req)
	default:
		return fmt.Errorf("radio request %s not supported", reqReq.Name)
	}
}

func insertRadioReq(rawbody []byte, w http.ResponseWriter, req *http.Request) error {
	paraReq := struct {
		Name     string `json:"name"`
		URI      string `json:"uri"`
		Descr    string `json:"descr"`
		PageIx   int    `json:"pageix"`
		PageSize int    `json:"pagesize"`
	}{}
	if err := json.Unmarshal(rawbody, &paraReq); err != nil {
		return err
	}
	log.Println("insert radio Request", paraReq)

	list := []*db.ResUriItem{}
	item := db.ResUriItem{
		URI:         paraReq.URI,
		Title:       paraReq.Name,
		Description: paraReq.Descr,
	}
	list = append(list, &item)

	trx, err := liteDB.GetNewTransaction()
	if err != nil {
		return err
	}

	err = liteDB.InsertRadioList(trx, list)
	if err != nil {
		return err
	}

	err = trx.Commit()
	if err != nil {
		return err
	}

	return fetchRadioReqInDB(paraReq.PageIx, paraReq.PageSize, w)
}

func fetchRadioReq(rawbody []byte, w http.ResponseWriter, req *http.Request) error {
	paraReq := struct {
		PageIx   int `json:"pageix"`
		PageSize int `json:"pagesize"`
	}{}
	if err := json.Unmarshal(rawbody, &paraReq); err != nil {
		return err
	}
	return fetchRadioReqInDB(paraReq.PageIx, paraReq.PageSize, w)
}

func fetchRadioReqInDB(pageIx, pageSize int, w http.ResponseWriter) error {
	log.Println("radio Request fetch ", pageIx, pageSize)
	list, err := liteDB.FetchRadio(pageIx, pageSize)
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
		PageIx: pageIx,
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
