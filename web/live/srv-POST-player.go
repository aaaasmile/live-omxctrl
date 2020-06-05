package live

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func handleChangeVolume(w http.ResponseWriter, req *http.Request) error {
	rawbody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	reqVol := struct {
		VolumeType string `json:"volume"`
	}{}

	if err := json.Unmarshal(rawbody, &reqVol); err != nil {
		return err
	}

	log.Println("Change volume request ", reqVol)

	pl := OmxPlayer{}
	switch reqVol.VolumeType {
	case "up":
		err = pl.VolumeUp()
	case "down":
		err = pl.VolumeDown()
	default:
		return fmt.Errorf("Change volume request not recognized ", reqVol)
	}
	return err
}

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

	return err
}

func handlePause(w http.ResponseWriter, req *http.Request) error {
	log.Println("Pause request ")
	pl := OmxPlayer{}
	return pl.Pause()
}
