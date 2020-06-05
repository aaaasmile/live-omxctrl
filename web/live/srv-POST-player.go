package live

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

func handleTogglePowerState(w http.ResponseWriter, req *http.Request) error {
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

	pl := OmxPlayer{}
	switch reqPower.PowerState {
	case "off":
		err = pl.PowerOff()
	case "on":
		err = startOmxPlayer("http://stream.srg-ssr.ch/m/rsc_de/aacp_96")

	default:
		return fmt.Errorf("Toggle power state  not recognized %s", reqPower.PowerState)
	}
	return err
}

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
	case "mute":
		err = pl.VolumeMute()
	case "unmute":
		err = pl.VolumeUnmute()
	default:
		return fmt.Errorf("Change volume request not recognized %s", reqVol.VolumeType)
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

func startOmxPlayer(URI string) error {
	log.Println("Start player wit URI ", URI)
	cmd := "omxplayer"
	args := []string{"-o", "local", URI}
	go func() {
		out, err := exec.Command(cmd, args...).Output()
		if err != nil {
			log.Printf("Error on executing omxplayer: %v", err)
		}
		log.Println("execute returns ", string(out))
	}()

	return nil
}
