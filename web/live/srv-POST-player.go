package live

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/aaaasmile/live-omxctrl/web/live/omx"
)

func handlePlayYoutube(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
	rawbody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	reqURI := struct {
		URI string `json:"uri"`
	}{}
	if err := json.Unmarshal(rawbody, &reqURI); err != nil {
		return err
	}
	log.Println("Play youtube URL ", reqURI.URI)

	if err := playYoutubeUri(reqURI.URI, pl); err != nil {
		return err
	}
	return returnStatus(w, req, pl)
}

func handleNextTitle(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
	if err := pl.NextTitle(); err != nil {
		return err
	}
	return returnStatus(w, req, pl)
}

func handlePreviousTitle(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
	if err := pl.PreviousTitle(); err != nil {
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
		return nil
	case "on":
		//u := "http://stream.srg-ssr.ch/m/rsc_de/aacp_96"
		//err = playOmxUri(u, pl)
		//u := "https://www.youtube.com/watch?v=3czUk1MmmvA"
		//err = playYoutubeUri(u, true, pl)
		err = pl.PowerOn()
	default:
		return fmt.Errorf("Toggle power state  not recognized %s", reqPower.PowerState)
	}
	if err != nil {
		return err
	}

	return returnStatusAfterCheck(w, req, pl)
}

func playOmxUri(u string, pl *omx.OmxPlayer) error {
	sleepTime := 400
	err := pl.StartOmxPlayer(u)
	return checkAfterStartPlay(sleepTime, err, pl)
}

func playYoutubeUri(u string, pl *omx.OmxPlayer) error {
	sleepTime := 700
	err := pl.StartYoutubeLink(u)
	return checkAfterStartPlay(sleepTime, err, pl)
}

func checkAfterStartPlay(sleepTime int, err error, pl *omx.OmxPlayer) error {
	time.Sleep(200 * time.Millisecond)
	i := 0
	for i < 8 {
		err = pl.CheckStatus()
		if err != nil {
			log.Println("Error and retry ", i, err)
			i++
		} else {
			break
		}
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
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
	return returnStatusAfterCheck(w, req, pl)
}

func returnStatusAfterCheck(w http.ResponseWriter, req *http.Request, pl *omx.OmxPlayer) error {
	res := struct {
		Player        string `json:"player"`
		Mute          string `json:"mute"`
		URI           string `json:"uri"`
		TrackDuration string `json:"trackDuration"`
		TrackPosition string `json:"trackPosition"`
		TrackStatus   string `json:"trackStatus"`
		Type          string `json:"type"`
	}{
		Player:        pl.GetStatePlaying(),
		Mute:          pl.GetStateMute(),
		URI:           pl.GetCurrURI(),
		TrackDuration: pl.TrackDuration,
		TrackPosition: pl.TrackPosition,
		TrackStatus:   pl.TrackStatus,
		Type:          "status",
	}

	return writeResponse(w, res)
}
