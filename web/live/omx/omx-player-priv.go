package omx

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/godbus/dbus"
)

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

func (op *OmxPlayer) getProperty(prop string) (*dbus.Variant, error) {
	op.mutex.Lock()
	defer op.mutex.Unlock()

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
	op.mutex.Lock()
	defer op.mutex.Unlock()

	if err := op.connectObjectDbBus(); err != nil {
		return err
	}
	op.coDBus.Call(action, 0)
	return nil
}

func (op *OmxPlayer) callIntAction(action string, id int) error {
	op.mutex.Lock()
	defer op.mutex.Unlock()

	if err := op.connectObjectDbBus(); err != nil {
		return err
	}
	op.coDBus.Call(action, 0, id)
	return nil
}

func (op *OmxPlayer) callStrAction(action string, para string) error {
	op.mutex.Lock()
	defer op.mutex.Unlock()

	if err := op.connectObjectDbBus(); err != nil {
		return err
	}
	op.coDBus.Call(action, 0, para)
	return nil
}

func (op *OmxPlayer) clearStatus() {
	op.TrackDuration = ""
	op.TrackPosition = ""
	op.TrackStatus = ""
}
