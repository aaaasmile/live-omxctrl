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

func testSomeDbus() error {
	log.Println("Let's try some dbus")
	//busAddr := "unix:abstract=/tmp/dbus-1OTLRLIFgE,guid=39be549b2196c379ccdf29585ed9674d"
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

func connectObjectDbBus() (dbus.BusObject, error) {
	u, err := user.Current()
	log.Println("User ", u.Username)
	fname := fmt.Sprintf("/tmp/omxplayerdbus.%s", u.Username)
	if _, err := os.Stat(fname); err == nil {
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
