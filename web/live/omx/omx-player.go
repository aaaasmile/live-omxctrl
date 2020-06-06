package omx

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"

	"github.com/godbus/dbus"
)

type OmxPlayer struct {
	coDBus  dbus.BusObject
	cmdOmx  *exec.Cmd
	currURI string
}

func NewOmxPlayer() *OmxPlayer {
	res := OmxPlayer{}
	return &res
}

func (op *OmxPlayer) StartOmxPlayer(URI string) error {
	if op.currURI == URI && op.cmdOmx != nil {
		log.Println("Same URI and player is active. Simple play")
		return op.callSimpleAction("Play")
	}
	if op.cmdOmx != nil {
		op.cmdOmx.Process.Kill()
	}
	log.Println("Start player wit URI ", URI)

	cmd := "omxplayer"
	args := []string{"-o", "local", URI}
	log.Println("turn on the player")
	op.cmdOmx = exec.Command(cmd, args...)
	if err := op.cmdOmx.Start(); err != nil {
		return fmt.Errorf("Error on executing omxplayer: %v", err)
	}
	op.currURI = URI
	return nil
}

func (op *OmxPlayer) getProperty(prop string) (*dbus.Variant, error) {
	if err := op.connectObjectDbBus(); err != nil {
		return nil, err
	}
	res, err := op.coDBus.GetProperty(prop)
	if err != nil {
		op.coDBus = nil
		return nil, err
	}
	return &res, nil
}

func (op *OmxPlayer) callSimpleAction(action string) error {
	if err := op.connectObjectDbBus(); err != nil {
		return err
	}
	op.coDBus.Call(action, 0)
	return nil
}

func (op *OmxPlayer) callIntAction(action string, id int) error {
	if err := op.connectObjectDbBus(); err != nil {
		return err
	}
	op.coDBus.Call(action, 0, id)
	return nil
}

func (op *OmxPlayer) CheckStatus() error {
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

	log.Println("Duration, position,  status ", dur, pos, status)
	return nil
}

func (op *OmxPlayer) Resume() error {
	log.Println("Resume")
	op.callSimpleAction("Play")
	return nil
}

func (op *OmxPlayer) Pause() error {
	log.Println("Pause")
	op.callSimpleAction("Pause")
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
	return nil
}

func (op *OmxPlayer) VolumeUnmute() error {
	log.Println("Volume Unmute")
	op.callSimpleAction("Unmute")
	return nil
}

func (op *OmxPlayer) PowerOff() error {
	log.Println("Power off, exit app")
	op.callIntAction("Action", 15)
	if op.cmdOmx != nil {
		op.cmdOmx.Process.Kill()
		op.cmdOmx = nil
		op.currURI = ""
	}
	return nil
}

func (op *OmxPlayer) connectObjectDbBus() error {
	if op.coDBus != nil {
		return nil
	}
	u, err := user.Current()
	log.Println("User ", u.Username)
	fname := fmt.Sprintf("/tmp/omxplayerdbus.%s", u.Username)
	if _, err := os.Stat(fname); err == nil {
		//busAddr := "unix:abstract=/tmp/dbus-1OTLRLIFgE,guid=39be549b2196c379ccdf29585ed9674d"
		raw, err := ioutil.ReadFile(fname)
		if err != nil {
			return err
		}
		os.Setenv("DBUS_SESSION_BUS_ADDRESS", string(raw))
		log.Println("Env DBUS_SESSION_BUS_ADDRESS set to ", string(raw))
	}

	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	obj := conn.Object("org.mpris.MediaPlayer2.omxplayer", "/org/mpris/MediaPlayer2/omxplayer")

	address := os.Getenv("DBUS_SESSION_BUS_ADDRESS")
	log.Println("session bus addr ", address)
	log.Println("dbus connection: ", conn)

	op.coDBus = obj
	return nil
}
