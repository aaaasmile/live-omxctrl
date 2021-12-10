package radio

import (
	"fmt"
	"log"
	"strings"

	"github.com/aaaasmile/live-omxctrl/conf"
	"github.com/aaaasmile/live-omxctrl/db"
	"github.com/aaaasmile/live-omxctrl/web/idl"
)

type infoFile struct {
	Title       string
	Description string
	Genre       string
}

type RadioPlayer struct {
	URI     string
	Info    *infoFile
	LiteDB  *db.LiteDB
	chClose chan struct{}
}

func (rp *RadioPlayer) IsUriForMe(uri string) bool {
	if strings.Contains(uri, "http") &&
		(strings.Contains(uri, "mp3") || strings.Contains(uri, "aacp") || strings.Contains(uri, "/live")) {
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

func (rp *RadioPlayer) SetURI(uri string) {
	rp.URI = uri
}

func (rp *RadioPlayer) GetPropValue(propname string) string {
	if propname == "genre" {
		return rp.Info.Genre
	}
	return ""
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
func (rp *RadioPlayer) GetStreamerCmd(cmdLineArr []string) (string, string) {
	args := strings.Join(cmdLineArr, " ")
	params := fmt.Sprintf("%s %s", args, rp.URI)
	return conf.Current.Player.Path, params
}
func (rp *RadioPlayer) CheckStatus(chDbOperation chan *idl.DbOperation) error {
	if rp.Info == nil {
		resItem, err := rp.LiteDB.FetchRadioFromURI(rp.URI)
		if err != nil {
			return err
		}
		info := infoFile{
			Title:       resItem.Title,
			Description: resItem.Description,
			Genre:       resItem.Genre,
		}
		log.Println("Radio info from db: ", resItem)
		hi := db.ResUriItem{
			URI:         rp.URI,
			Title:       info.Title,
			Description: info.Genre,
			Type:        rp.Name(),
		}
		dop := idl.DbOperation{
			DbOpType: idl.DbOpHistoryInsert,
			Payload:  hi,
		}
		chDbOperation <- &dop
		rp.Info = &info
	}

	return nil
}

func (rp *RadioPlayer) CreateStopChannel() chan struct{} {
	if rp.chClose == nil {
		rp.chClose = make(chan struct{})
	}
	return rp.chClose
}

func (rp *RadioPlayer) GetCmdStopChannel() chan struct{} {
	return rp.chClose
}

func (rp *RadioPlayer) CloseStopChannel() {
	if rp.chClose != nil {
		close(rp.chClose)
		rp.chClose = nil
	}
}

func (rp *RadioPlayer) GetTrackDuration() (string, bool) {
	return "", false
}
func (rp *RadioPlayer) GetTrackPosition() (string, bool) {
	return "", false
}
func (rp *RadioPlayer) GetTrackStatus() (string, bool) {
	return "", false
}
