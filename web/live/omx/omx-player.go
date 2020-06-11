package omx

import (
	"fmt"
	"log"
	"os/exec"
	"sync"

	"github.com/godbus/dbus"
)

type SPstateplaying int

const (
	SPoff = iota
	SPplaying
	SPpause
)

type SMstatemute int

const (
	SMnormal = iota
	SMmuted
)

type StateOmx struct {
	CurrURI      string
	StatePlaying SPstateplaying
	StateMute    SMstatemute
}

type OmxPlayer struct {
	coDBus        dbus.BusObject
	cmdOmx        *exec.Cmd
	mutex         *sync.Mutex
	State         StateOmx
	TrackDuration string
	TrackPosition string
	TrackStatus   string
}

func NewOmxPlayer() *OmxPlayer {
	res := OmxPlayer{
		mutex: &sync.Mutex{},
	}
	return &res
}

func (op *OmxPlayer) StartOmxPlayer(URI string, chst chan *StateOmx) error {
	op.mutex.Lock()
	defer op.mutex.Unlock()

	if op.State.CurrURI == URI && op.cmdOmx != nil {
		log.Println("Same URI and player is active. Simple play")
		return op.callSimpleAction("Play")
	}
	if op.cmdOmx != nil {
		op.cmdOmx.Process.Kill()
	}
	log.Println("Start player wit URI ", URI)

	cmd := "omxplayer"
	args := []string{"-o", "local", URI}
	op.cmdOmx = exec.Command(cmd, args...)
	op.execCommand(URI, chst)

	return nil
}

func (op *OmxPlayer) StartYoutubeLink(URI string, chst chan *StateOmx) error {
	op.mutex.Lock()
	defer op.mutex.Unlock()

	if op.cmdOmx != nil {
		op.cmdOmx.Process.Kill()
	}
	log.Println("Start youtube player wit URI ", URI)

	cmd := fmt.Sprintf("omxplayer -o local  `youtube-dl -f mp4 -g %s`", URI)
	op.cmdOmx = exec.Command("bash", "-c", cmd)
	op.execCommand(URI, chst)

	return nil
}

func (op *OmxPlayer) NextTitle(chst chan *StateOmx) error {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	if op.cmdOmx == nil {
		return nil
	}

	// Some try fo a playlist
	// Next Title works only if the the player is still running
	// Next title as radio swiss doesn't works. Better only with mp3.
	// Radio url as next, better to restart the player
	//u := "/home/igors/music/youtube/milanoda_bere_spot.mp3"
	u := "/home/igors/music/youtube/Elisa - Tua Per Sempre-3czUk1MmmvA.mp3"
	if op.State.CurrURI == u {
		// switch to test how to make a play list
		//u = "http://stream.srg-ssr.ch/m/rsc_de/aacp_96"
		u = "/home/igors/music/youtube/Gianna Nannini - Fenomenale (Official Video)-HKwWcJCtwck.mp3"
		//u = "https://www.youtube.com/watch?v=3czUk1MmmvA"
		//u = "`youtube-dl -f mp4 -g https://www.youtube.com/watch?v=3czUk1MmmvA`"
		//return op.StartOmxPlayer(u)
	}
	log.Println("Play the next title", u)
	op.callStrAction("OpenUri", u)

	chst <- &StateOmx{CurrURI: u, StatePlaying: SPplaying}
	return nil
}

func (op *OmxPlayer) CheckStatus() error {
	op.mutex.Lock()
	defer op.mutex.Unlock()

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

func (op *OmxPlayer) Resume(chst chan *StateOmx) error {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	if op.cmdOmx == nil {
		return nil
	}

	log.Println("Resume")
	op.callSimpleAction("Play")
	if op.State.CurrURI != "" {
		chst <- &StateOmx{CurrURI: op.State.CurrURI, StatePlaying: SPplaying}
	}

	return nil
}

func (op *OmxPlayer) Pause(chst chan *StateOmx) error {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	if op.cmdOmx == nil {
		return nil
	}

	log.Println("Pause")
	op.callSimpleAction("Pause")
	if op.State.CurrURI != "" {
		chst <- &StateOmx{CurrURI: op.State.CurrURI, StatePlaying: SPpause}
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

func (op *OmxPlayer) VolumeMute(chst chan *StateOmx) error {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	if op.cmdOmx == nil {
		return nil
	}

	log.Println("Volume Mute")
	op.callSimpleAction("Mute")
	if op.State.CurrURI != "" {
		chst <- &StateOmx{CurrURI: op.State.CurrURI, StateMute: SMmuted}
	}
	return nil
}

func (op *OmxPlayer) VolumeUnmute(chst chan *StateOmx) error {
	op.mutex.Lock()
	defer op.mutex.Unlock()
	if op.cmdOmx == nil {
		return nil
	}

	log.Println("Volume Unmute")
	op.callSimpleAction("Unmute")
	if op.State.CurrURI != "" {
		chst <- &StateOmx{CurrURI: op.State.CurrURI, StateMute: SMnormal}
	}

	return nil
}

func (op *OmxPlayer) PowerOff(chst chan *StateOmx) error {
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
	chst <- &StateOmx{CurrURI: "", StatePlaying: SPoff}
	op.coDBus = nil

	return nil
}

func (op *OmxPlayer) SetState(st *StateOmx) {
	log.Println("Set OmxPlayer state ", st)
	op.mutex.Lock()
	defer op.mutex.Unlock()
	op.State = *st
	if st.StatePlaying == SPoff {
		op.coDBus = nil
		op.clearTrackStatus()
	}
}
