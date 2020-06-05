package live

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/godbus/dbus"
)

func testSomeDbus() error {
	log.Println("Let's try some dbus")
	//busAddr := "unix:abstract=/tmp/dbus-1OTLRLIFgE,guid=39be549b2196c379ccdf29585ed9674d"
	//opath := dbus.ObjectPath(busAddr)
	u, err := user.Current()
	log.Println("User ", u.Username)
	fname := fmt.Sprintf("/tmp/omxplayerdbus.%s", u.Username)
	if _, err := os.Stat(fname); err == nil {
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
	address := os.Getenv("DBUS_SESSION_BUS_ADDRESS")
	log.Println("session bus addr ", address)
	log.Println("dbus connection: ", conn)

	// var s []string
	// err = conn.Object.Call("/org/mpris/MediaPlayer2")
	// //.Ge org.freedesktop.DBus.Properties.Get", 0).Store(&s)
	// if err != nil {
	// 	return fmt.Errorf("Failed to get list of owned names: %v", err)
	// }
	// log.Println("Currently owned names on the session bus:")
	// for _, v := range s {
	// 	fmt.Println(v)
	// }

	return nil
}

func testSomeDbus1() error {
	log.Println("Let's try some dbus")
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	log.Println("dbus connection: ", conn)
	defer conn.Close()

	var s []string
	err = conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&s)
	if err != nil {
		return fmt.Errorf("Failed to get list of owned names: %v", err)
	}
	log.Println("Currently owned names on the session bus:")
	for _, v := range s {
		fmt.Println(v)
	}

	return nil
}
