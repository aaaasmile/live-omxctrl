package radio

import (
	"fmt"
	"log"
	"strings"

	"github.com/aaaasmile/live-omxctrl/db"
)

type infoFile struct {
	Title       string
	Description string
}

type RadioPlayer struct {
	URI     string
	Info    *infoFile
	chClose chan struct{}
}

func (rp *RadioPlayer) IsUriForMe(uri string) bool {
	if strings.Contains(uri, "http") &&
		(strings.Contains(uri, "mp3") || strings.Contains(uri, "aacp")) {
		log.Println("This is a streaming resource ", uri)
		rp.URI = uri
		return true
	}
	return false
}

func (rp *RadioPlayer) GetStatusSleepTime() int {
	return 500
}

func (rp *RadioPlayer) GetURI() string {
	return rp.URI
}
func (rp *RadioPlayer) GetTitle() string {
	if rp.Info != nil {
		return rp.Info.Title
	}
	return ""
}
func (rp *RadioPlayer) GetDescription() string {
	if rp.Info != nil {
		return rp.Info.Description
	}
	return ""
}
func (rp *RadioPlayer) Name() string {
	return "radio"
}
func (rp *RadioPlayer) GetStreamerCmd(cmdLineArr []string) string {
	args := strings.Join(cmdLineArr, " ")
	cmd := fmt.Sprintf("omxplayer %s %s", args, rp.URI)
	return cmd
}
func (rp *RadioPlayer) CheckStatus(chHistoryItem chan *db.HistoryItem) (bool, error) {
	if rp.Info == nil {
		info := infoFile{
			// TODO read from db afetr file scan
		}
		hi := db.HistoryItem{
			URI:         rp.URI,
			Title:       info.Title,
			Description: info.Description,
			Type:        rp.Name(),
		}
		chHistoryItem <- &hi
	}
	return false, nil
}

func (rp *RadioPlayer) GetStopChannel() chan struct{} {
	if rp.chClose == nil {
		rp.chClose = make(chan struct{})
	}
	return rp.chClose
}

func (rp *RadioPlayer) CloseStopChannel() {
	if rp.chClose != nil {
		close(rp.chClose)
		rp.chClose = nil
	}
}
