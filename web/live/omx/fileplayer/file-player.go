package fileplayer

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aaaasmile/live-omxctrl/db"
)

type infoFile struct {
	Title         string
	Description   string
	DurationInSec int
}

type FilePlayer struct {
	URI     string
	Info    *infoFile
	chClose chan struct{}
}

func (fp *FilePlayer) IsUriForMe(uri string) bool {
	if strings.Contains(uri, "/home") &&
		(strings.Contains(uri, ".mp3") || strings.Contains(uri, ".ogg") || strings.Contains(uri, ".wav")) {
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
func (fp *FilePlayer) CheckStatus(chHistoryItem chan *db.HistoryItem) (bool, error) {
	if fp.Info == nil {
		info := infoFile{
			// TODO read from db afetr file scan
		}
		hi := db.HistoryItem{
			URI:           fp.URI,
			Title:         info.Title,
			Description:   info.Description,
			DurationInSec: info.DurationInSec,
			Type:          fp.Name(),
			Duration:      time.Duration(int64(info.DurationInSec) * int64(time.Second)).String(),
		}
		chHistoryItem <- &hi
	}
	return false, nil
}

func (fp *FilePlayer) GetStopChannel() chan struct{} {
	if fp.chClose == nil {
		fp.chClose = make(chan struct{})
	}
	return fp.chClose
}

func (fp *FilePlayer) CloseStopChannel() {
	if fp.chClose != nil {
		close(fp.chClose)
		fp.chClose = nil
	}
}
