package omx

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"

	"github.com/aaaasmile/live-omxctrl/web/idl"
	"github.com/aaaasmile/live-omxctrl/web/live/omx/playlist"
	"github.com/godbus/dbus"
)

type actionTD int

///home/igors/projects/go/bin/stringer -type=actionTD

const (
	actTerminate actionTD = iota
	actPlaying
	actPause
)

type actionDef struct {
	URI    string
	Action actionTD
}

func (op *OmxPlayer) execCommand(uri, cmdText string, chstop chan struct{}) {
	log.Println("Prepare to start the player with execCommand")
	go func(cmdText string, actCh chan *actionDef, uri string, chstop chan struct{}) {
		// op.cmdOmx = exec.Command("bash", "-c", cmd)
		log.Println("Submit the command ", cmdText)
		//cmd := exec.Command("omxplayer", "-o", "hdmi", "/home/igors/Music/youtube/gianna-fenomenale.mp3")
		//cmd := exec.Command("bash", "-c", "omxplayer -o luz /home/igors/Music/youtube/gianna-fenomenale.mp3")
		cmd := exec.Command("bash", "-c", cmdText)
		actCh <- &actionDef{
			URI:    uri,
			Action: actPlaying,
		}
		// out, err := cmdOmx.Output()
		// log.Println("Command out ", string(out))
		// if err != nil {
		// 	log.Println("Command executed with error: ", err)
		// }
		var stdoutBuf bytes.Buffer
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)

		if err := cmd.Start(); err == nil {
			done := make(chan error, 1)
			go func() {
				done <- cmd.Wait()
			}()
		loop:
			for {

				select {
				case <-chstop:
					log.Println("Received stop signal")
					if err := cmd.Process.Kill(); err != nil {
						log.Println("Error on killing the process ", err)
						break loop
					}
				case err := <-done:
					log.Println("Process finished")
					if err != nil {
						log.Println("Error on process termination")
					}
					log.Println(string(stdoutBuf.Bytes()))
					break loop
				default:
					outStr := string(stdoutBuf.Bytes())
					if outStr != "" {
						log.Println("***> ", outStr)
					}
				}
			}
			log.Println("Exit from loop")

		} else {
			log.Println("ERROR cmd.Run() failed with", err)
		}

		log.Println("Player has been terminated. Cmd was ", cmdText)
		actCh <- &actionDef{
			URI:    uri,
			Action: actTerminate,
		}

	}(cmdText, op.chAction, uri, chstop)
}

func listenStateAction(actCh chan *actionDef, op *OmxPlayer) {
	log.Println("Waiting for action to change the state")
	var stateCurrent SPstateplaying
	stateCurrent = SPoff
	uriPlaying := ""
	for {
		st := <-actCh
		log.Println("New action in state: ", st.Action.String(), stateCurrent.String())
		stateNext := StateOmx{CurrURI: st.URI, StatePlayer: SPundef}
		switch stateCurrent {
		case SPoff:
			switch st.Action {
			case actPlaying:
				stateNext.StatePlayer = SPplaying
				uriPlaying = st.URI
			}
		case SPplaying:
			switch st.Action {
			case actPlaying:
				stateNext.StatePlayer = SPrestart
				uriPlaying = st.URI
			case actPause:
				stateNext.StatePlayer = SPpause
			case actTerminate:
				stateNext.StatePlayer = SPoff
				uriPlaying = ""
			}
		case SPpause:
			switch st.Action {
			case actPlaying:
				stateNext.StatePlayer = SPplaying
			}
		case SPrestart:
			switch st.Action {
			case actTerminate:
				stateNext.StatePlayer = SPplaying
			}
		}

		log.Println("Calculated next state ", stateNext.StatePlayer.String())
		if stateNext.StatePlayer != SPundef {
			log.Println("State trigger a change")
			stateCurrent = stateNext.StatePlayer
			stateNext.CurrURI = uriPlaying
			op.mutex.Lock()
			op.setState(&stateNext)
			op.mutex.Unlock()
		} else {
			log.Println("Ignored action ", st.Action.String())
		}
	}
}

func (op *OmxPlayer) startPlayListCurrent(prov idl.StreamProvider) error {
	log.Println("Start current item ", op.PlayList)
	var curr *playlist.PlayItem
	var ok bool
	if curr, ok = op.PlayList.CheckCurrent(); !ok {
		return nil
	}
	log.Println("Current item is ", curr)
	op.mutex.Lock()
	defer op.mutex.Unlock()

	if op.state.CurrURI != "" {
		log.Println("Shutting down the current player of ", op.state.CurrURI)
		//op.callIntAction("Action", 15)
		if pp, ok := op.Providers[op.state.CurrURI]; ok {
			chStop := pp.GetStopChannel()
			chStop <- struct{}{}
			pp.CloseStopChannel()
			op.Providers[op.state.CurrURI] = nil
		}
	}
	uri := prov.GetURI()
	op.Providers[uri] = prov

	log.Println("Start player wit URI ", uri)

	if len(op.cmdLineArr) == 0 {
		return fmt.Errorf("Command line is not set")
	}
	cmd := prov.GetStreamerCmd(op.cmdLineArr)
	log.Println("Start the command: ", cmd)
	//op.cmdOmx = exec.Command("bash", "-c", cmd)
	op.execCommand(uri, cmd, prov.GetStopChannel())
	//op.setState(&StateOmx{CurrURI: uri, StatePlaying: SPplaying})

	return nil
}

func (op *OmxPlayer) connectObjectDbBus() error {
	if op.coDBus != nil {
		return nil
	}
	u, err := user.Current()
	log.Println("User ", u.Username)

	// TODO multiple istance of omxplayerdbus are possible
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

func (op *OmxPlayer) callStrAction(action string, para string) error {
	if err := op.connectObjectDbBus(); err != nil {
		return err
	}
	op.coDBus.Call(action, 0, para)
	return nil
}

func (op *OmxPlayer) clearTrackStatus() {
	op.TrackDuration = ""
	op.TrackPosition = ""
	op.TrackStatus = ""
}

func (op *OmxPlayer) setState(st *StateOmx) {
	log.Println("Set OmxPlayer state ", st)
	op.state.CurrURI = st.CurrURI
	op.state.StateMute = st.StateMute
	op.state.StatePlayer = st.StatePlayer
	op.state.Info = st.Info
	if st.StatePlayer == SPoff {
		op.coDBus = nil
		op.clearTrackStatus()
	}
	op.chstatus <- &op.state
}
