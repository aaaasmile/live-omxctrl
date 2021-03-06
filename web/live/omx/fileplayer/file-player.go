package fileplayer

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aaaasmile/live-omxctrl/db"
	"github.com/aaaasmile/live-omxctrl/web/idl"
	"github.com/aaaasmile/live-omxctrl/web/live/omx/dbus"
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

type FilePlayer struct {
	URI     string
	Info    *infoFile
	Dbus    *dbus.OmxDbus
	chClose chan struct{}
}

func (fp *FilePlayer) IsUriForMe(uri string) bool {
	if strings.Contains(uri, "/home") &&
		(strings.Contains(uri, ".mp4") || strings.Contains(uri, ".avi") || strings.Contains(uri, ".mkv") ||
			strings.Contains(uri, ".mp3") || strings.Contains(uri, ".ogg") || strings.Contains(uri, ".wav")) {
		log.Println("this is a music file ", uri)
		fp.URI = uri
		return true
	}
	return false
}

func (fp *FilePlayer) GetStatusSleepTime() int {
	return 300
}

func (fp *FilePlayer) GetURI() string {
	return fp.URI
}
func (fp *FilePlayer) GetTitle() string {
	if fp.Info != nil {
		return fp.Info.Title
	}
	return ""
}
func (fp *FilePlayer) GetDescription() string {
	if fp.Info != nil {
		return fp.Info.Description
	}
	return ""
}
func (fp *FilePlayer) Name() string {
	return "file"
}
func (fp *FilePlayer) GetStreamerCmd(cmdLineArr []string) string {
	args := strings.Join(cmdLineArr, " ")
	cmd := fmt.Sprintf("omxplayer %s %s", args, fp.URI)
	return cmd
}
func (fp *FilePlayer) CheckStatus(chDbOperation chan *idl.DbOperation) error {
	st := &omxstate.StateOmx{}
	if err := fp.Dbus.CheckTrackStatus(st); err != nil {
		return err
	}

	if fp.Info == nil {
		info := infoFile{
			// TODO read from db
		}
		info.DurationInSec, _ = strconv.Atoi(st.TrackDuration)
		info.TrackDuration = time.Duration(int64(info.DurationInSec) * int64(time.Second)).String()
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

func (fp *FilePlayer) CreateStopChannel() chan struct{} {
	if fp.chClose == nil {
		fp.chClose = make(chan struct{})
	}
	return fp.chClose
}

func (fp *FilePlayer) GetCmdStopChannel() chan struct{} {
	return fp.chClose
}

func (fp *FilePlayer) CloseStopChannel() {
	if fp.chClose != nil {
		close(fp.chClose)
		fp.chClose = nil
	}
}

func (fp *FilePlayer) GetTrackDuration() (string, bool) {
	if fp.Info != nil {
		return fp.Info.TrackDuration, true
	}
	return "", false

}
func (fp *FilePlayer) GetTrackPosition() (string, bool) {
	if fp.Info != nil {
		return fp.Info.TrackPosition, true
	}
	return "", false

}
func (fp *FilePlayer) GetTrackStatus() (string, bool) {
	if fp.Info != nil {
		return fp.Info.TrackStatus, true
	}
	return "", false
}
