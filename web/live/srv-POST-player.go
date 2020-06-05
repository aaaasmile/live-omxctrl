package live

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func handlePlay(w http.ResponseWriter, req *http.Request) error {
	rawbody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	reqPlay := struct {
		URI string
	}{}

	if err := json.Unmarshal(rawbody, &reqPlay); err != nil {
		return err
	}

	log.Println("Play request ", reqPlay)

	pl := OmxPlayer{}
	if reqPlay.URI != "" {
		err = pl.OpenUri(reqPlay.URI)
	} else {
		err = pl.Resume()
	}

	if err != nil {
		return err
	}

	return nil
}

func handlePause(w http.ResponseWriter, req *http.Request) error {
	log.Println("Pause request ")
	pl := OmxPlayer{}
	return pl.Pause()
}
