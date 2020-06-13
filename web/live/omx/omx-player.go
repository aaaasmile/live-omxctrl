package omx

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"

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
		op.cmdOmx.Process.Kill()
	}
	log.Println("Start youtube player wit URI ", URI)

	cmd := fmt.Sprintf("omxplayer -o local  `youtube-dl -f mp4 -g %s`", URI)
	op.cmdOmx = exec.Command("bash", "-c", cmd)
	op.execCommand()
	op.setState(&StateOmx{CurrURI: URI, StatePlaying: SPplaying})

	return nil
}

func (op *OmxPlayer) NextTitle() error {
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
	if op.state.CurrURI == u {
		// switch to test how to make a play list
		//u = "http://stream.srg-ssr.ch/m/rsc_de/aacp_96"
		u = "/home/igors/music/youtube/Gianna Nannini - Fenomenale (Official Video)-HKwWcJCtwck.mp3"
		//u = "https://www.youtube.com/watch?v=3czUk1MmmvA"
		//u = "`youtube-dl -f mp4 -g https://www.youtube.com/watch?v=3czUk1MmmvA`"
		//return op.StartOmxPlayer(u)
	}
	log.Println("Play the next title", u)
	op.callStrAction("OpenUri", u)

	op.setState(&StateOmx{CurrURI: u, StatePlaying: SPplaying})
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
		op.setState(&StateOmx{CurrURI: op.state.CurrURI, StateMute: SMmuted})
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
		op.setState(&StateOmx{CurrURI: op.state.CurrURI, StateMute: SMnormal})
	}

	return nil
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
		op.clearTrackStatus()
	}
	op.chstatus <- &op.state
}
