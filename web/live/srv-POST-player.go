package live

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aaaasmile/live-omxctrl/web/live/omx"
)

func handleTogglePowerState(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
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

	log.Println("Toggle power state request ", reqPower)

	switch reqPower.PowerState {
	case "off":
		err = pl.PowerOff()
	case "on":
		err = pl.StartOmxPlayer("http://stream.srg-ssr.ch/m/rsc_de/aacp_96")

	default:
		return fmt.Errorf("Toggle power state  not recognized %s", reqPower.PowerState)
	}
	return err
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
	return err
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

	return err
}

func handlePause(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
	log.Println("Pause request ")
	return pl.Pause()
}

func handlePlayerState(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
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
