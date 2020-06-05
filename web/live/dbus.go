package live

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/godbus/dbus"
)

var sdbus *dbus.Conn

func testSomeDbus() error {
	log.Println("Let's try some dbus")
	//busAddr := "unix:abstract=/tmp/dbus-1OTLRLIFgE,guid=39be549b2196c379ccdf29585ed9674d"
	//opath := dbus.ObjectPath(busAddr)
	var err error
	if sdbus == nil {
		sdbus, err = connectDbBus()
		if err != nil {
			return err
		}
	}

	obj := sdbus.Object("org.mpris.MediaPlayer2.omxplayer", "/org/mpris/MediaPlayer2/omxplayer")
	log.Println("Object is ", obj)
	dur, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Duration")
	if err != nil {
		return err
	}
	log.Println("Duration is ", dur)
	return nil
}

func connectDbBus() (*dbus.Conn, error) {
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
	address := os.Getenv("DBUS_SESSION_BUS_ADDRESS")
	log.Println("session bus addr ", address)
	log.Println("dbus connection: ", conn)

	return conn, nil
}

// func testSomeDbus1() error {
// 	log.Println("Let's try some dbus")
// 	conn, err := dbus.SessionBus()
// 	if err != nil {
// 		return err
// 	}
// 	log.Println("dbus connection: ", conn)
// 	defer conn.Close()

// 	var s []string
// 	err = conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&s)
// 	if err != nil {
// 		return fmt.Errorf("Failed to get list of owned names: %v", err)
// 	}
// 	log.Println("Currently owned names on the session bus:")
// 	for _, v := range s {
// 		fmt.Println(v)
// 	}

// 	return nil
// }
