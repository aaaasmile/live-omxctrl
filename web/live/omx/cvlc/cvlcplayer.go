package cvlc

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aaaasmile/live-omxctrl/db"
	"github.com/aaaasmile/live-omxctrl/web/idl"
	"github.com/aaaasmile/live-omxctrl/web/live/omx/omxstate"
)

type infoFile struct {
	Title         string
	Description   string
	DurationInSec int
	TrackDuration string
	TrackPosition string
	TrackStatus   string
}

type CvlcPlayer struct {
	URI     string
	Info    *infoFile
	LiteDB  *db.LiteDB
	chClose chan struct{}
}

func (fp *CvlcPlayer) IsUriForMe(uri string) bool {
	if strings.Contains(uri, "/home") &&
		(strings.Contains(uri, ".mkv") ||
			strings.Contains(uri, ".avi")) {
		log.Println("this is a mkv file ", uri)
		fp.URI = uri
		return true
	}
	return false
}

func (fp *CvlcPlayer) GetStatusSleepTime() int {
	return 300
}

func (fp *CvlcPlayer) GetURI() string {
	return fp.URI
}

func (rp *CvlcPlayer) SetURI(uri string) {
	rp.URI = uri
}

func (fp *CvlcPlayer) GetTitle() string {
	if fp.Info != nil {
		return fp.Info.Title
	}
	return ""
}

func (rp *CvlcPlayer) GetPropValue(propname string) string {
	return ""
}

func (fp *CvlcPlayer) GetDescription() string {
	if fp.Info != nil {
		return fp.Info.Description
	}
	return ""
}
func (fp *CvlcPlayer) Name() string {
	return "file"
}
func (fp *CvlcPlayer) GetStreamerCmd(cmdLineArr []string) (string, string, []string) {
	moreargs := []string{}
	params := fmt.Sprintf("%s %s", "--aout=alsa --alsa-audio-device=plughw:b1,0", fp.URI)
	return "cvlc", params, moreargs
}
func (fp *CvlcPlayer) CheckStatus(chDbOperation chan *idl.DbOperation) error {
	st := &omxstate.StateOmx{}

	if fp.Info == nil {
		titles := strings.Split(fp.URI, "/")
		title := ""
		if len(titles) > 0 {
			title = titles[len(titles)-1]
		}
		info := infoFile{
			// TODO read from db using URI -> see player radio
			Title: title,
		}
		info.DurationInSec, _ = strconv.Atoi(st.TrackDuration)
		if info.DurationInSec > 0 {
			info.TrackDuration = time.Duration(int64(info.DurationInSec) * int64(time.Second)).String()
		}

		hi := db.ResUriItem{
			URI:           fp.URI,
			Title:         info.Title,
			Description:   info.Description,
			DurationInSec: info.DurationInSec,
			Type:          fp.Name(),
			Duration:      info.TrackDuration,
		}
		dop := idl.DbOperation{
			DbOpType: idl.DbOpHistoryInsert,
			Payload:  hi,
		}
		chDbOperation <- &dop
		fp.Info = &info
		log.Println("file-player info status set")
	}

	fp.Info.TrackPosition = st.TrackPosition
	fp.Info.TrackStatus = st.TrackStatus
	log.Println("Status set to ", fp.Info)
	return nil
}

func (fp *CvlcPlayer) CreateStopChannel() chan struct{} {
	if fp.chClose == nil {
		fp.chClose = make(chan struct{})
	}
	return fp.chClose
}

func (fp *CvlcPlayer) GetCmdStopChannel() chan struct{} {
	return fp.chClose
}

func (fp *CvlcPlayer) CloseStopChannel() {
	if fp.chClose != nil {
		close(fp.chClose)
		fp.chClose = nil
	}
}

func (fp *CvlcPlayer) GetTrackDuration() (string, bool) {
	if fp.Info != nil {
		return fp.Info.TrackDuration, true
	}
	return "", false

}
func (fp *CvlcPlayer) GetTrackPosition() (string, bool) {
	if fp.Info != nil {
		return fp.Info.TrackPosition, true
	}
	return "", false

}
func (fp *CvlcPlayer) GetTrackStatus() (string, bool) {
	if fp.Info != nil {
		return fp.Info.TrackStatus, true
	}
	return "", false
}
