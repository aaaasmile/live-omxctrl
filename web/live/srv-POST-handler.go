package live

import (
	"fmt"
	"log"
	"net/http"
)

func handlePost(w http.ResponseWriter, req *http.Request) error {
	var err error
	lastPath := getURLForRoute(req.RequestURI)
	log.Println("Check the last path ", lastPath)
	switch lastPath {
	case "PlayURI": // TODO provide status on all handlers
		err = handlePlay(w, req, player)
	case "Pause":
		err = handlePause(w, req, player)
	case "ChangeVolume":
		err = handleChangeVolume(w, req, player)
	case "SetPowerState":
		err = handleSetPowerState(w, req, player)
	case "GetPlayerState":
		err = handlePlayerState(w, req, player)
	case "NextTitle":
		err = handleNextTitle(w, req, player)
	default:
		return fmt.Errorf("%s method is not supported", lastPath)
	}

	return err
}
