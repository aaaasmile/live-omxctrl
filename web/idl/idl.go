package idl

import "github.com/aaaasmile/live-omxctrl/db"

var (
	Appname = "live-omxctrl"
	Buildnr = "00.01.22.20201227-00"
)

type StreamProvider interface {
	IsUriForMe(uri string) bool
	GetStatusSleepTime() int
	GetURI() string
	GetTitle() string
	GetDescription() string
	Name() string
	GetStreamerCmd(cmdLineArr []string) string
	CheckStatus(chHistoryItem chan *db.HistoryItem) error
	CreateStopChannel() chan struct{}
	GetCmdStopChannel() chan struct{}
	CloseStopChannel()
	GetTrackDuration() (string, bool)
	GetTrackPosition() (string, bool)
	GetTrackStatus() (string, bool)
}
