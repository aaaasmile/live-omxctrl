package omx

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/aaaasmile/live-omxctrl/db"
	"github.com/aaaasmile/live-omxctrl/web/idl"
	"github.com/aaaasmile/live-omxctrl/web/live/omx/omxstate"
	"github.com/aaaasmile/live-omxctrl/web/live/omx/playlist"
	"github.com/godbus/dbus"
)

type OmxPlayer struct {
	coDBus        dbus.BusObject
	mutex         *sync.Mutex
	state         omxstate.StateOmx
	chHistoryItem chan *db.HistoryItem
	TrackDuration string
	TrackPosition string
	TrackStatus   string
	cmdLineArr    []string
	PlayList      *playlist.LLPlayList
	Providers     map[string]idl.StreamProvider
	ChAction      chan *omxstate.ActionDef
	ChStatus      chan *omxstate.StateOmx
}

func NewOmxPlayer(chhisitem chan *db.HistoryItem) *OmxPlayer {
	cha := make(chan *omxstate.ActionDef)
	chst := make(chan *omxstate.StateOmx)
	res := OmxPlayer{
		mutex:         &sync.Mutex{},
		chHistoryItem: chhisitem,
		cmdLineArr:    make([]string, 0),
		Providers:     make(map[string]idl.StreamProvider),
		ChAction:      cha,
		ChStatus:      chst,
	}

	go res.listenStatus(chst)

	return &res
}

func (op *OmxPlayer) SetCommandLine(commaline string) {
	op.cmdLineArr = make([]string, 0)
	arr := strings.Split(commaline, ",")
	for _, item := range arr {
		if len(item) > 0 {
			op.cmdLineArr = append(op.cmdLineArr, item)
		}
	}
	log.Println("Command line set to ", commaline, op.cmdLineArr)
}

func (op *OmxPlayer) GetStatePlaying() string {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	return op.state.StatePlayer.String()
}

func (op *OmxPlayer) GetStateMute() string {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	return op.state.StateMute.String()
}

func (op *OmxPlayer) GetStateTitle() string {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	if prov, ok := op.Providers[op.state.CurrURI]; ok {
		return prov.GetTitle()
	}

	return ""
}

func (op *OmxPlayer) GetStateDescription() string {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	if prov, ok := op.Providers[op.state.CurrURI]; ok {
		return prov.GetDescription()
	}

	return ""
}

func (op *OmxPlayer) GetCurrURI() string {
	// please do not call this after a mutex lock
	log.Println("getCurrURI")
	op.mutex.Lock()
	defer op.mutex.Unlock()
	return op.state.CurrURI
}

func (op *OmxPlayer) StartPlay(URI string, prov idl.StreamProvider) error {
	var err error
	if op.PlayList, err = playlist.CreatePlaylistFromProvider(URI, prov); err != nil {
		return err
	}
	log.Println("StartPlay ", URI)

	return op.startPlayListCurrent(prov)
}

func (op *OmxPlayer) PreviousTitle() (string, error) {
	if op.PlayList == nil {
		log.Println("Nothing to play because no playlist is provided")
		return "", nil
	}
	var curr *playlist.PlayItem
	var ok bool
	if _, ok = op.PlayList.CheckCurrent(); !ok {
		return "", nil
	}

	op.mutex.Lock()
	defer op.mutex.Unlock()

	if op.state.CurrURI == "" {
		log.Println("Player is not active, ignore next title")
		return "", nil
	}

	if curr, ok = op.PlayList.Previous(); !ok {
		return "", nil
	}

	u := curr.URI
	log.Println("the previous title is", u)

	return u, nil
}

func (op *OmxPlayer) NextTitle() (string, error) {
	if op.PlayList == nil {
		log.Println("Nothing to play because no playlist is provided")
		return "", nil
	}
	var curr *playlist.PlayItem
	var ok bool
	if _, ok = op.PlayList.CheckCurrent(); !ok {
		return "", nil
	}

	op.mutex.Lock()
	defer op.mutex.Unlock()

	if op.state.CurrURI == "" {
		return "", nil
	}

	if curr, ok = op.PlayList.Next(); !ok {
		return "", nil
	}

	u := curr.URI
	log.Println("the next title is", u)

	return u, nil
}

func (op *OmxPlayer) CheckStatus(uri string) error {
	if uri == "" {
		return nil
	}
	op.mutex.Lock()
	defer op.mutex.Unlock()

	log.Println("Check status req", op.state)
	op.clearTrackStatus()

	if prov, ok := op.Providers[op.state.CurrURI]; ok {
		completed, err := prov.CheckStatus(op.chHistoryItem)
		if err != nil {
			return err
		}
		if completed {
			log.Println("Check status completed")
			return nil
		}
	}
	log.Println("go ahead with dbus status")
	dur, err := op.getProperty("org.mpris.MediaPlayer2.Player.Duration")
	if err != nil {
		return err
	}
	pos, err := op.getProperty("org.mpris.MediaPlayer2.Player.Position")
	if err != nil {
		return err
	}

	status, err := op.getProperty("org.mpris.MediaPlayer2.Player.PlaybackStatus")
	if err != nil {
		return err
	}

	op.TrackDuration = fmt.Sprint(dur)
	op.TrackPosition = fmt.Sprint(pos)
	op.TrackStatus = fmt.Sprint(status)

	log.Println("Duration, position,  status ", dur, pos, status)
	return nil
}

func (op *OmxPlayer) Resume() error {
	op.mutex.Lock()
	defer op.mutex.Unlock()

	if op.state.CurrURI != "" {
		log.Println("Resume")
		op.callSimpleAction("Play")
		op.ChAction <- &omxstate.ActionDef{Action: omxstate.ActPlaying}
	}

	return nil
}

func (op *OmxPlayer) Pause() error {
	op.mutex.Lock()
	defer op.mutex.Unlock()

	if op.state.CurrURI != "" {
		log.Println("Pause")
		op.callSimpleAction("Pause")
		op.ChAction <- &omxstate.ActionDef{Action: omxstate.ActPause}
	}
	return nil
}

func (op *OmxPlayer) VolumeUp() error {
	op.mutex.Lock()
	defer op.mutex.Unlock()

	if op.state.CurrURI != "" {
		log.Println("VolumeUp")
		// dbus-send --print-reply=literal --session --dest=org.mpris.MediaPlayer2.omxplayer /org/mpris/MediaPlayer2 org.mpris.MediaPlayer2.Player.Action int32:18 >/dev/null
		op.callIntAction("Action", 18)
	}
	// TODO check the volume level
	return nil
}

func (op *OmxPlayer) VolumeDown() error {
	op.mutex.Lock()
	defer op.mutex.Unlock()

	if op.state.CurrURI != "" {
		log.Println("VolumeDown")
		op.callIntAction("Action", 17)
	}
	// TODO check the volume level
	return nil
}

func (op *OmxPlayer) VolumeMute(chStateRsp chan *omxstate.StateOmx) error {
	op.mutex.Lock()
	defer op.mutex.Unlock()

	if (op.state.StatePlayer == omxstate.SPplaying) &&
		(op.state.StateMute == omxstate.SMnormal) {
		log.Println("Volume Mute")
		op.callSimpleAction("Mute")

		op.ChAction <- &omxstate.ActionDef{Action: omxstate.ActMute, ChStateRsp: chStateRsp}
	} else {
		log.Println("Ignore Mute request in state ", op.state)
		chStateRsp <- &op.state
	}
	return nil
}

func (op *OmxPlayer) VolumeUnmute(chStateRsp chan *omxstate.StateOmx) error {
	op.mutex.Lock()
	defer op.mutex.Unlock()

	if (op.state.StatePlayer == omxstate.SPplaying) &&
		(op.state.StateMute == omxstate.SMmuted) {
		log.Println("Volume Unmute")
		op.callSimpleAction("Unmute")
		op.ChAction <- &omxstate.ActionDef{Action: omxstate.ActUnmute, ChStateRsp: chStateRsp}
	} else {
		log.Println("Ignore Unmute request in state ", op.state)
		chStateRsp <- &op.state
	}

	return nil
}

func (op *OmxPlayer) PowerOff() error {
	op.mutex.Lock()
	defer op.mutex.Unlock()

	log.Println("Power off, terminate omxplayer with signal kill")
	op.freeAllProviders()
	return nil
}

func (op *OmxPlayer) freeAllProviders() {
	for k, prov := range op.Providers {
		log.Println("Sending kill signal to ", k)
		ch := prov.GetStopChannel()
		ch <- struct{}{}
		prov.CloseStopChannel()
	}

	op.Providers = make(map[string]idl.StreamProvider)

}
