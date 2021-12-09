package omx

import (
	"fmt"
	"log"

	"github.com/aaaasmile/live-omxctrl/web/idl"
	"github.com/aaaasmile/live-omxctrl/web/live/omx/playlist"
)

func (op *OmxPlayer) startPlayListCurrent(prov idl.StreamProvider) error {
	log.Println("Start current item ", op.PlayList)
	var curr *playlist.PlayItem
	var ok bool
	if curr, ok = op.PlayList.CheckCurrent(); !ok {
		return nil
	}
	log.Println("Current item is ", curr)
	op.mutex.Lock()
	defer op.mutex.Unlock()

	curURI := op.state.CurrURI
	if curURI != "" {
		log.Println("Shutting down the current player of ", curURI)
		if pp, ok := op.Providers[curURI]; ok {
			chStop := pp.GetCmdStopChannel()
			if chStop != nil {
				chStop <- struct{}{}
				pp.CloseStopChannel()
			}
			delete(op.Providers, curURI)
		}
	}
	uri := prov.GetURI()
	if uri == "" {
		return fmt.Errorf("URI is not recognized in player")
	}
	op.Providers[uri] = prov

	log.Println("Start player with URI ", uri)

	if len(op.cmdLineArr) == 0 {
		return fmt.Errorf("Command line is not set")
	}
	cmd := prov.GetStreamerCmd(op.cmdLineArr)
	log.Println("Start the command: ", cmd)
	op.execCommand(uri, cmd, prov.CreateStopChannel())

	return nil
}
