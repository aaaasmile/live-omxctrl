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
	log.Println("recieved readio request ", reqReq)
	switch reqReq.Name {
	case "FetchRadio":
		return fetchRadioReq(rawbody, w, req)
	case "Insert":
		return insertRadioReq(rawbody, w, req)
	case "Edit":
		return editRadioReq(rawbody, w, req)
	case "Delete":
		return deleteRadioReq(rawbody, w, req)
	default:
		return fmt.Errorf("radio request %s not supported", reqReq.Name)
	}
}

func deleteRadioReq(rawbody []byte, w http.ResponseWriter, req *http.Request) error {
	paraReq := struct {
		ID       int `json:"id"`
		PageIx   int `json:"pageix"`
		PageSize int `json:"pagesize"`
	}{}
	if err := json.Unmarshal(rawbody, &paraReq); err != nil {
		return err
	}
	log.Println("delete radio Request", paraReq)

	item := db.ResUriItem{
		ID: paraReq.ID,
	}

	trx, err := liteDB.GetNewTransaction()
	if err != nil {
		return err
	}

	err = liteDB.DeleteRadioItem(trx, item)
	if err != nil {
		return err
	}

	err = trx.Commit()
	if err != nil {
		return err
	}

	return fetchRadioReqInDB(paraReq.PageIx, paraReq.PageSize, w)
}

func editRadioReq(rawbody []byte, w http.ResponseWriter, req *http.Request) error {
	paraReq := struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		URI         string `json:"uri"`
		Description string `json:"description"`
		PageIx      int    `json:"pageix"`
		PageSize    int    `json:"pagesize"`
	}{}
	if err := json.Unmarshal(rawbody, &paraReq); err != nil {
		return err
	}
	log.Println("edit radio Request", paraReq)
	if paraReq.URI == "" || paraReq.Title == "" {
		return fmt.Errorf("mandatory fields are required")
	}

	item := db.ResUriItem{
		ID:          paraReq.ID,
		URI:         paraReq.URI,
		Title:       paraReq.Title,
		Description: paraReq.Description,
	}

	trx, err := liteDB.GetNewTransaction()
	if err != nil {
		return err
	}

	err = liteDB.EditRadioItem(trx, item)
	if err != nil {
		return err
	}

	err = trx.Commit()
	if err != nil {
		return err
	}

	return fetchRadioReqInDB(paraReq.PageIx, paraReq.PageSize, w)
}

func insertRadioReq(rawbody []byte, w http.ResponseWriter, req *http.Request) error {
	paraReq := struct {
		Title       string `json:"title"`
		URI         string `json:"uri"`
		Description string `json:"description"`
		PageIx      int    `json:"pageix"`
		PageSize    int    `json:"pagesize"`
	}{}
	if err := json.Unmarshal(rawbody, &paraReq); err != nil {
		return err
	}
	log.Println("insert radio Request", paraReq)
	if paraReq.URI == "" || paraReq.Title == "" {
		return fmt.Errorf("mandatory fields are required")
	}

	list := []*db.ResUriItem{}
	item := db.ResUriItem{
		URI:         paraReq.URI,
		Title:       paraReq.Title,
		Description: paraReq.Description,
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
