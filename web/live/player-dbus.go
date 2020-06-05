package live

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/godbus/dbus"
)

var coDBus dbus.BusObject

type OmxPlayer struct {
}

func (op *OmxPlayer) CheckStatus() error {
	var err error
	if coDBus == nil {
		coDBus, err = connectObjectDbBus()
		if err != nil {
			return err
		}
	}

	log.Println("Object is ", coDBus)
	dur, err := coDBus.GetProperty("org.mpris.MediaPlayer2.Player.Duration")
	if err != nil {
		coDBus = nil
		return err
	}
	pos, err := coDBus.GetProperty("org.mpris.MediaPlayer2.Player.Position")
	if err != nil {
		coDBus = nil
		return err
	}

	status, err := coDBus.GetProperty("org.mpris.MediaPlayer2.Player.PlaybackStatus")
	if err != nil {
		coDBus = nil
		return err
	}

	log.Println("Duration, position,  status ", dur, pos, status)
	return nil
}

func (op *OmxPlayer) OpenUri(uri string) error {
	var err error
	if coDBus == nil {
		coDBus, err = connectObjectDbBus()
		if err != nil {
			return err
		}
	}

	coDBus.Call("OpenUri", 0, uri)
	log.Println("open uri: ", uri)

	return nil
}

func (op *OmxPlayer) Resume() error {
	var err error
	if coDBus == nil {
		coDBus, err = connectObjectDbBus()
		if err != nil {
			return err
		}
	}
	// coDBus.Call("Action", 0, 16) + https://github.com/popcornmix/omxplayer/blob/master/KeyConfig.h
	coDBus.Call("Play", 0)
	log.Println("Play")
	return nil
}

func (op *OmxPlayer) Pause() error {
	var err error
	if coDBus == nil {
		coDBus, err = connectObjectDbBus()
		if err != nil {
			return err
		}
	}
	// coDBus.Call("Action", 0, 16) +
	coDBus.Call("Pause", 0)
	log.Println("Pause")
	return nil
}

func (op *OmxPlayer) VolumeUp() error {
	var err error
	if coDBus == nil {
		coDBus, err = connectObjectDbBus()
		if err != nil {
			return err
		}
	}
	// dbus-send --print-reply=literal --session --dest=org.mpris.MediaPlayer2.omxplayer /org/mpris/MediaPlayer2 org.mpris.MediaPlayer2.Player.Action int32:18 >/dev/null
	coDBus.Call("Action", 0, 18)
	log.Println("VolumeUp") // ACTION_INCREASE_VOLUME = 18,
	return nil
}

func (op *OmxPlayer) VolumeDown() error {
	var err error
	if coDBus == nil {
		coDBus, err = connectObjectDbBus()
		if err != nil {
			return err
		}
	}
	coDBus.Call("Action", 0, 17) // ACTION_DECREASE_VOLUME = 17,
	log.Println("VolumeDown")
	return nil
}

func (op *OmxPlayer) VolumeMute() error {
	var err error
	if coDBus == nil {
		coDBus, err = connectObjectDbBus()
		if err != nil {
			return err
		}
	}
	coDBus.Call("Mute", 0)
	log.Println("Volume Mute")
	return nil
}

func (op *OmxPlayer) VolumeUnmute() error {
	var err error
	if coDBus == nil {
		coDBus, err = connectObjectDbBus()
		if err != nil {
			return err
		}
	}
	coDBus.Call("Unmute", 0)
	log.Println("Volume Unmute")
	return nil
}

func (op *OmxPlayer) PowerOff() error {
	var err error
	if coDBus == nil {
		coDBus, err = connectObjectDbBus()
		if err != nil {
			return err
		}
	}
	coDBus.Call("Action", 0, 15) // ACTION_EXIT = 15,
	log.Println("Power off, exit app")
	return nil
}

func connectObjectDbBus() (dbus.BusObject, error) {
	u, err := user.Current()
	log.Println("User ", u.Username)
	fname := fmt.Sprintf("/tmp/omxplayerdbus.%s", u.Username)
	if _, err := os.Stat(fname); err == nil {
		//busAddr := "unix:abstract=/tmp/dbus-1OTLRLIFgE,guid=39be549b2196c379ccdf29585ed9674d"
		raw, err := ioutil.ReadFile(fname)
		if err != nil {
			return nil, err
		}
		os.Setenv("DBUS_SESSION_BUS_ADDRESS", string(raw))
		log.Println("Env DBUS_SESSION_BUS_ADDRESS set to ", string(raw))
	}

	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}
	obj := conn.Object("org.mpris.MediaPlayer2.omxplayer", "/org/mpris/MediaPlayer2/omxplayer")

	address := os.Getenv("DBUS_SESSION_BUS_ADDRESS")
	log.Println("session bus addr ", address)
	log.Println("dbus connection: ", conn)

	return obj, nil
}
