package live

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aaaasmile/live-omxctrl/web/live/omx"
)

func handleNextTitle(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
	if err := pl.NextTitle(); err != nil {
		return err
	}
	return returnStatus(w, req, pl)
}

func handleSetPowerState(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
	rawbody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	reqPower := struct {
		PowerState string `json:"power"`
	}{}

	if err := json.Unmarshal(rawbody, &reqPower); err != nil {
		return err
	}

	log.Println("Set power state request ", reqPower)

	switch reqPower.PowerState {
	case "off":
		err = pl.PowerOff()
	case "on":
		//u := "http://stream.srg-ssr.ch/m/rsc_de/aacp_96"
		//u := "`youtube-dl -f mp4 -g https://www.youtube.com/watch?v=3czUk1MmmvA`"
		u := "https://www.youtube.com/watch?v=3czUk1MmmvA"
		//err = pl.StartOmxPlayer(u)
		//time.Sleep(800 * time.Millisecond)
		err = pl.StartYoutubeLink(u)
		return nil

	default:
		return fmt.Errorf("Toggle power state  not recognized %s", reqPower.PowerState)
	}

	return returnStatus(w, req, pl)
}

func handleChangeVolume(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
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

	switch reqVol.VolumeType {
	case "up":
		err = pl.VolumeUp()
	case "down":
		err = pl.VolumeDown()
	case "mute":
		err = pl.VolumeMute()
	case "unmute":
		err = pl.VolumeUnmute()
	default:
		return fmt.Errorf("Change volume request not recognized %s", reqVol.VolumeType)
	}

	return returnStatus(w, req, pl)
}

func handlePlay(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
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

	if reqPlay.URI != "" {
		err = pl.StartOmxPlayer(reqPlay.URI)
	} else {
		err = pl.Resume()
	}

	return returnStatus(w, req, pl)
}

func handlePause(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
	log.Println("Pause request ")
	pl.Pause()

	return returnStatus(w, req, pl)
}

func handlePlayerState(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
	return returnStatus(w, req, pl)
}

func returnStatus(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
	if err := pl.CheckStatus(); err != nil {
		return err
	}

	res := struct {
		Player        string `json:"player"`
		Mute          string `json:"mute"`
		URI           string `json:"uri"`
		TrackDuration string `json:"trackDuration"`
		TrackPosition string `json:"trackPosition"`
		TrackStatus   string `json:"trackStatus"`
	}{
		Player:        pl.StatePlaying,
		Mute:          pl.StateMute,
		URI:           pl.CurrURI,
		TrackDuration: pl.TrackDuration,
		TrackPosition: pl.TrackPosition,
		TrackStatus:   pl.TrackStatus,
	}

	return writeResponse(w, res)
}
