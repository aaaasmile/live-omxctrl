//+build windows

package dbus

import (
	"github.com/aaaasmile/live-omxctrl/web/live/omx/omxstate"
	"github.com/godbus/dbus"
)

type OmxDbus struct {
}

func (op *OmxDbus) ClearDbus() {
}

func (op *OmxDbus) connectObjectDbBus() error {
	return nil
}

func (op *OmxDbus) getProperty(prop string) (*dbus.Variant, error) {
	return nil, nil
}

func (op *OmxDbus) CallSimpleAction(action string) error {
	return nil
}

func (op *OmxDbus) CallIntAction(action string, id int) error {
	return nil
}

func (op *OmxDbus) CallStrAction(action string, para string) error {
	return nil
}

func (op *OmxDbus) CheckTrackStatus(st *omxstate.StateOmx) error {

	return nil
}
