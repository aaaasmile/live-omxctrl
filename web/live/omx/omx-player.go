package omx

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"

	"github.com/aaaasmile/live-omxctrl/db"
	"github.com/aaaasmile/live-omxctrl/web/idl"
	"github.com/aaaasmile/live-omxctrl/web/live/omx/playlist"
	"github.com/godbus/dbus"
)

type OmxPlayer struct {
	coDBus        dbus.BusObject
	cmdOmx        *exec.Cmd
	mutex         *sync.Mutex
	state         StateOmx
	chstatus      chan *StateOmx
	chHistoryItem chan *db.HistoryItem
	TrackDuration string
	TrackPosition string
	TrackStatus   string
	cmdLineArr    []string
	PlayList      *playlist.LLPlayList
	Providers     map[string]idl.StreamProvider
	chAction      chan *actionDef
}

func NewOmxPlayer(chst chan *StateOmx, chhisitem chan *db.HistoryItem) *OmxPlayer {
	cha := make(chan *actionDef)
	res := OmxPlayer{
		mutex:         &sync.Mutex{},
		chstatus:      chst,
		chHistoryItem: chhisitem,
		cmdLineArr:    make([]string, 0),
		Providers:     make(map[string]idl.StreamProvider),
		chAction:      cha,
	}
	go listenStateAction(res.chAction, &res)

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
	log.Println("In Mx GetStatePlaying ")
	return op.state.StatePlayer.String()
}

func (op *OmxPlayer) GetStateMute() string {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	log.Println("In Mx GetStateMute ")
	return op.state.StateMute.String()
}

func (op *OmxPlayer) GetStateTitle() string {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	log.Println("In Mx GetStateTitle ")
	if prov, ok := op.Providers[op.state.CurrURI]; ok {
		return prov.GetTitle()
	}

	return ""
}

func (op *OmxPlayer) GetStateDescription() string {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	log.Println("In Mx GetStateDescription ")
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
	log.Println("In Mx GetCurrURI ")
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
	log.Println("In Mx PreviousTitle")

	if op.cmdOmx == nil {
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
	log.Println("In Mx NextTitle")

	if op.cmdOmx == nil {
		log.Println("Player is not active, ignore next title")
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
	op.mutex.Lock()
	defer op.mutex.Unlock()
	log.Println("In Mx CheckStatus")

	log.Println("Check status req", op.state)
	op.clearTrackStatus()
	if op.cmdOmx == nil {
		return nil
	}
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
	log.Println("In Mx Resume")
	if op.cmdOmx == nil {
		return nil
	}

	log.Println("Resume")
	op.callSimpleAction("Play")
	op.chAction <- &actionDef{Action: actPlaying}

	return nil
}

func (op *OmxPlayer) Pause() error {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	log.Println("In Mx Pause")
	if op.cmdOmx == nil {
		return nil
	}

	log.Println("Pause")
	op.callSimpleAction("Pause")
	op.chAction <- &actionDef{Action: actPause}
	return nil
}

func (op *OmxPlayer) VolumeUp() error {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	log.Println("In Mx VolumeUp")
	if op.cmdOmx == nil {
		return nil
	}

	log.Println("VolumeUp")
	// dbus-send --print-reply=literal --session --dest=org.mpris.MediaPlayer2.omxplayer /org/mpris/MediaPlayer2 org.mpris.MediaPlayer2.Player.Action int32:18 >/dev/null
	op.callIntAction("Action", 18)
	// ACTION_INCREASE_VOLUME = 18,
	// TODO check the volume level
	return nil
}

func (op *OmxPlayer) VolumeDown() error {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	log.Println("In Mx VolumeDown")
	if op.cmdOmx == nil {
		return nil
	}

	log.Println("VolumeDown")
	op.callIntAction("Action", 17) // ACTION_DECREASE_VOLUME = 17,
	// TODO check the volume level
	return nil
}

func (op *OmxPlayer) VolumeMute() error {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	log.Println("In Mx VolumeMute")
	if op.cmdOmx == nil {
		return nil
	}

	log.Println("Volume Mute")
	op.callSimpleAction("Mute")
	if op.state.CurrURI != "" {
		op.setState(&StateOmx{StatePlayer: op.state.StatePlayer, CurrURI: op.state.CurrURI, StateMute: SMmuted})
	}
	return nil
}

func (op *OmxPlayer) VolumeUnmute() error {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	log.Println("In Mx VolumeUnmute")
	if op.cmdOmx == nil {
		return nil
	}

	log.Println("Volume Unmute")
	op.callSimpleAction("Unmute")
	if op.state.CurrURI != "" {
		op.setState(&StateOmx{StatePlayer: op.state.StatePlayer, CurrURI: op.state.CurrURI, StateMute: SMnormal})
	}

	return nil
}

func (op *OmxPlayer) PowerOff() error {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	log.Println("In Mx PowerOff")
	if op.cmdOmx == nil {
		log.Println("Player is not active. Nothing to do")
		return nil
	}

	log.Println("Power off, terminate omxplayer")
	op.callIntAction("Action", 15)
	if op.cmdOmx != nil {
		op.cmdOmx.Process.Kill()
		op.cmdOmx = nil

	}
	op.setState(&StateOmx{StatePlayer: SPoff})

	return nil
}
