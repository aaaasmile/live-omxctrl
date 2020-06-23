package omx

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/aaaasmile/live-omxctrl/playlist"
	"github.com/godbus/dbus"
)

type OmxPlayer struct {
	coDBus        dbus.BusObject
	cmdOmx        *exec.Cmd
	mutex         *sync.Mutex
	state         StateOmx
	chstatus      chan *StateOmx
	TrackDuration string
	TrackPosition string
	TrackStatus   string
	cmdLine       []string
	startedTime   time.Time
	PlayList      *playlist.LLPlayList
}

func NewOmxPlayer(chst chan *StateOmx) *OmxPlayer {
	res := OmxPlayer{
		mutex:    &sync.Mutex{},
		chstatus: chst,
	}
	return &res
}

func (op *OmxPlayer) SetCommandLine(commaline string) {
	op.cmdLine = make([]string, 0)
	arr := strings.Split(commaline, ",")
	for _, item := range arr {
		if len(item) > 0 {
			op.cmdLine = append(op.cmdLine, item)
		}
	}
}

func (op *OmxPlayer) GetStatePlaying() string {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	return op.state.StatePlaying.String()
}

func (op *OmxPlayer) GetStateMute() string {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	return op.state.StateMute.String()
}

func (op *OmxPlayer) GetCurrURI() string {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	return op.state.CurrURI
}

func (op *OmxPlayer) StartOmxPlayer(URI string) error {
	op.mutex.Lock()
	defer op.mutex.Unlock()

	if op.state.CurrURI == URI && op.cmdOmx != nil {
		log.Println("Same URI and player is active. Simple play")
		return op.callSimpleAction("Play")
	}
	if op.cmdOmx != nil {
		op.cmdOmx.Process.Kill()
	}
	log.Println("Start player wit URI ", URI)

	cmd := "omxplayer"
	args := append(op.cmdLine, URI)
	log.Println("Command line is: ", cmd, args)
	op.cmdOmx = exec.Command(cmd, args...)
	op.execCommand()
	op.setState(&StateOmx{CurrURI: URI, StatePlaying: SPplaying})

	return nil
}

func (op *OmxPlayer) StartYoutubeLink(URI string) error {
	op.mutex.Lock()
	defer op.mutex.Unlock()

	if op.cmdOmx != nil {
		log.Println("Shuttind down the current player", op.cmdOmx)
		op.callIntAction("Action", 15)
		op.cmdOmx.Process.Kill()
	}
	log.Println("Start youtube player wit URI ", URI)

	args := strings.Join(op.cmdLine, " ")
	cmd := fmt.Sprintf("omxplayer %s `youtube-dl -f mp4 -g %s`", args, URI)
	op.cmdOmx = exec.Command("bash", "-c", cmd)
	op.execCommand()
	op.setState(&StateOmx{CurrURI: URI, StatePlaying: SPplaying})

	return nil
}

func (op *OmxPlayer) PreviousTitle() error {
	// TODO...
	return fmt.Errorf("TODO...")
	// if op.PlayList.FirstItem == nil {
	// 	return nil
	// }
	// op.mutex.Lock()
	// defer op.mutex.Unlock()
	// if op.cmdOmx == nil {
	// 	return nil
	// }

	// // TODO check the start time and if it is small then play op.PlayList.CurrItem.Previous
	// u := op.state.CurrURI

	// log.Println("Play the previous title", u)
	// op.callStrAction("OpenUri", u)

	// op.setState(&StateOmx{CurrURI: u, StatePlaying: SPplaying})
	// return nil
}

func (op *OmxPlayer) NextTitle() error {
	if op.PlayList == nil {
		log.Println("Nothing to play because no playlist is provided")
		return nil
	}
	var curr *playlist.PlayItem
	var ok bool
	if _, ok = op.PlayList.CheckCurrent(); !ok {
		return nil
	}

	op.mutex.Lock()

	if op.cmdOmx == nil {
		op.mutex.Unlock()
		log.Println("Player is not active, ignore next title")
		return nil
	}

	if curr, ok = op.PlayList.Next(); !ok {
		op.mutex.Unlock()
		return nil
	}

	if curr.ItemType != playlist.PITMp3 {
		op.mutex.Unlock()
		return op.startCurrentItem()
	}
	u := curr.URI
	log.Println("Play the next title with action", u)
	op.callStrAction("OpenUri", u)

	op.setState(&StateOmx{CurrURI: u, StatePlaying: SPplaying})

	op.mutex.Unlock()
	return nil
}

func (op *OmxPlayer) CheckStatus() error {
	op.mutex.Lock()
	defer op.mutex.Unlock()

	log.Println("Check status req")
	op.clearTrackStatus()
	if op.cmdOmx == nil {
		return nil
	}
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
	if op.cmdOmx == nil {
		return nil
	}

	log.Println("Resume")
	op.callSimpleAction("Play")
	if op.state.CurrURI != "" {
		op.setState(&StateOmx{CurrURI: op.state.CurrURI, StatePlaying: SPplaying})
	}

	return nil
}

func (op *OmxPlayer) Pause() error {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	if op.cmdOmx == nil {
		return nil
	}

	log.Println("Pause")
	op.callSimpleAction("Pause")
	if op.state.CurrURI != "" {
		op.setState(&StateOmx{CurrURI: op.state.CurrURI, StatePlaying: SPpause})
	}
	return nil
}

func (op *OmxPlayer) VolumeUp() error {
	op.mutex.Lock()
	defer op.mutex.Unlock()
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
	if op.cmdOmx == nil {
		return nil
	}

	log.Println("Volume Mute")
	op.callSimpleAction("Mute")
	if op.state.CurrURI != "" {
		op.setState(&StateOmx{StatePlaying: op.state.StatePlaying, CurrURI: op.state.CurrURI, StateMute: SMmuted})
	}
	return nil
}

func (op *OmxPlayer) VolumeUnmute() error {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	if op.cmdOmx == nil {
		return nil
	}

	log.Println("Volume Unmute")
	op.callSimpleAction("Unmute")
	if op.state.CurrURI != "" {
		op.setState(&StateOmx{StatePlaying: op.state.StatePlaying, CurrURI: op.state.CurrURI, StateMute: SMnormal})
	}

	return nil
}

func (op *OmxPlayer) PowerOn() error {
	log.Println("Powern on")
	var err error
	if op.PlayList, err = playlist.GetCurrentPlaylist(); err != nil {
		return err
	}
	log.Println("Start the play list ", op.PlayList.Name)

	return op.startCurrentItem()
}

func (op *OmxPlayer) PowerOff() error {
	op.mutex.Lock()
	defer op.mutex.Unlock()
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
	op.setState(&StateOmx{StatePlaying: SPoff})

	return nil
}

func (op *OmxPlayer) setState(st *StateOmx) {
	log.Println("Set OmxPlayer state ", st)
	op.state.CurrURI = st.CurrURI
	op.state.StateMute = st.StateMute
	op.state.StatePlaying = st.StatePlaying
	if st.StatePlaying == SPoff {
		op.coDBus = nil
		op.cmdOmx = nil
		op.clearTrackStatus()
	}
	op.chstatus <- &op.state
}
