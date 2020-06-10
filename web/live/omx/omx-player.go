package omx

import (
	"fmt"
	"log"
	"os/exec"
	"sync"

	"github.com/godbus/dbus"
)

type OmxPlayer struct {
	coDBus        dbus.BusObject
	cmdOmx        *exec.Cmd
	mutex         *sync.Mutex
	CurrURI       string
	StatePlaying  string
	StateMute     string
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

func (op *OmxPlayer) StartOmxPlayer(URI string) error {
	if op.CurrURI == URI && op.cmdOmx != nil {
		log.Println("Same URI and player is active. Simple play")
		return op.callSimpleAction("Play")
	}
	if op.cmdOmx != nil {
		op.cmdOmx.Process.Kill()
	}
	log.Println("Start player wit URI ", URI)

	op.mutex.Lock()
	defer op.mutex.Unlock()

	cmd := "omxplayer"
	args := []string{"-o", "local", URI}
	log.Println("turn on the player")
	//out, err := exec.Command(cmd, args...).Output()
	//log.Println("Out is ", string(out))
	//if err != nil {
	//	return fmt.Errorf("Error on executing omxplayer: %v", err)
	//}
	op.cmdOmx = exec.Command(cmd, args...)
	go func() {
		//out, err := exec.Command("bash", "-c", cmd).Output()
		out, err := op.cmdOmx.Output()
		if err != nil {
			log.Println("Failed to execute command: ", err)
		}
		log.Println("Command out ", string(out))
		// TODO: this is only one function and the status should be set via channel(uri,playing,volume)
	}()

	// if err := op.cmdOmx.Start(); err != nil {
	// 	return fmt.Errorf("Error on executing omxplayer: %v", err)
	// }
	op.CurrURI = URI
	op.StatePlaying = "playing"
	return nil
}

func (op *OmxPlayer) StartYoutubeLink(URI string) error {
	// doesn't works
	if op.cmdOmx != nil {
		op.cmdOmx.Process.Kill()
	}
	log.Println("Start youtube player wit URI ", URI)

	op.mutex.Lock()
	defer op.mutex.Unlock()

	cmd := fmt.Sprintf("omxplayer -o local  `youtube-dl -f mp4 -g %s`", URI)
	op.cmdOmx = exec.Command("bash", "-c", cmd)

	go func() {
		//out, err := exec.Command("bash", "-c", cmd).Output()
		out, err := op.cmdOmx.Output()
		if err != nil {
			log.Println("Failed to execute command: ", err)
		}
		log.Println("Command out ", string(out))
	}()

	// op.cmdOmx = exec.Command("bash", "-c", cmd)
	// if err := op.cmdOmx.Start(); err != nil {
	// 	return fmt.Errorf("Error on executing omxplayer: %v", err)
	// }

	op.CurrURI = URI
	op.StatePlaying = "playing"
	return nil
}

func (op *OmxPlayer) NextTitle() error {
	// Some try fo a playlist
	// Next Title works only if the the player is still running
	// Next title as radio swiss doesn't works. Better only with mp3.
	// Radio url as next, better to restart the player
	//u := "/home/igors/music/youtube/milanoda_bere_spot.mp3"
	u := "/home/igors/music/youtube/Elisa - Tua Per Sempre-3czUk1MmmvA.mp3"
	if op.CurrURI == u {
		// switch to test how to make a play list
		//u = "http://stream.srg-ssr.ch/m/rsc_de/aacp_96"
		u = "/home/igors/music/youtube/Gianna Nannini - Fenomenale (Official Video)-HKwWcJCtwck.mp3"
		//u = "https://www.youtube.com/watch?v=3czUk1MmmvA"
		//u = "`youtube-dl -f mp4 -g https://www.youtube.com/watch?v=3czUk1MmmvA`"
		//return op.StartOmxPlayer(u)
	}
	log.Println("Play the next title", u)
	op.callStrAction("OpenUri", u)

	op.CurrURI = u
	op.StatePlaying = "playing"
	return nil
}

func (op *OmxPlayer) CheckStatus() error {
	op.clearStatus()
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
	log.Println("Resume")
	op.callSimpleAction("Play")
	if op.CurrURI != "" {
		op.StatePlaying = "playing"
	}

	return nil
}

func (op *OmxPlayer) Pause() error {
	log.Println("Pause")
	op.callSimpleAction("Pause")
	if op.CurrURI != "" {
		op.StatePlaying = "pause"
	}

	return nil
}

func (op *OmxPlayer) VolumeUp() error {
	log.Println("VolumeUp")
	// dbus-send --print-reply=literal --session --dest=org.mpris.MediaPlayer2.omxplayer /org/mpris/MediaPlayer2 org.mpris.MediaPlayer2.Player.Action int32:18 >/dev/null
	op.callIntAction("Action", 18)
	// ACTION_INCREASE_VOLUME = 18,
	return nil
}

func (op *OmxPlayer) VolumeDown() error {
	log.Println("VolumeDown")
	op.callIntAction("Action", 17) // ACTION_DECREASE_VOLUME = 17,
	return nil
}

func (op *OmxPlayer) VolumeMute() error {
	log.Println("Volume Mute")
	op.callSimpleAction("Mute")
	if op.CurrURI != "" {
		op.StateMute = "muted"
	}
	return nil
}

func (op *OmxPlayer) VolumeUnmute() error {
	log.Println("Volume Unmute")
	op.callSimpleAction("Unmute")
	if op.CurrURI != "" {
		op.StateMute = ""
	}

	return nil
}

func (op *OmxPlayer) PowerOff() error {
	if op.cmdOmx == nil {
		log.Println("Player is not active. Nothing to do")
	}
	log.Println("Power off, exit app")
	op.callIntAction("Action", 15)
	if op.cmdOmx != nil {
		op.cmdOmx.Process.Kill()
		op.cmdOmx = nil

	}
	op.CurrURI = ""
	op.StatePlaying = ""
	op.StateMute = ""
	op.coDBus = nil
	return nil
}
