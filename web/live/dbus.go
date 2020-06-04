package live

import (
	"fmt"
	"log"

	"github.com/godbus/dbus"
)

func testSomeDbus() error {
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
