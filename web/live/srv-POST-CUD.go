package live

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func handleCUD(w http.ResponseWriter, req *http.Request) error {
	rawbody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	reqReq := struct {
		Table  string `json:"table"`
		Method string `json:"method"`
	}{}
	if err := json.Unmarshal(rawbody, &reqReq); err != nil {
		return err
	}
	if reqReq.Table == "" || reqReq.Method == "" {
		return fmt.Errorf("malformed request (table or method)")
	}

	switch reqReq.Table {
	case "Radio":
		switch reqReq.Method {
		case "insert":
			return insertRadioDbReq(rawbody, w, req)
		}
	}

	return fmt.Errorf("CUD request not supported: %s - %s", reqReq.Table, reqReq.Method)
}
