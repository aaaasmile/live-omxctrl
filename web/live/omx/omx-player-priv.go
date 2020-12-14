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
	"syscall"

	"github.com/aaaasmile/live-omxctrl/web/idl"
	"github.com/aaaasmile/live-omxctrl/web/live/omx/omxstate"
	"github.com/aaaasmile/live-omxctrl/web/live/omx/playlist"
	"github.com/godbus/dbus"
)

func (op *OmxPlayer) execCommand(uri, cmdText string, chstop chan struct{}) {
	log.Println("Prepare to start the player with execCommand")
	go func(cmdText string, actCh chan *omxstate.ActionDef, uri string, chstop chan struct{}) {
		log.Println("Submit the command in background ", cmdText)
		cmd := exec.Command("bash", "-c", cmdText)
		actCh <- &omxstate.ActionDef{
			URI:    uri,
			Action: omxstate.ActPlaying,
		}

		var stdoutBuf, stderrBuf bytes.Buffer
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
		cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

		if err := cmd.Start(); err == nil {
			log.Println("PID started ", cmd.Process.Pid)
			done := make(chan error, 1)
			go func() {
				done <- cmd.Wait()
				log.Println("Wait ist terminated")
			}()

			select {
			case <-chstop:
				log.Println("Received stop signal, kill parent and child processes")
				if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
					log.Println("Error on killing the process ", err)
				}
			case err := <-done:
				log.Println("Process finished")
				if err != nil {
					log.Println("Error on process termination =>", err)
				}
				log.Println(string(stderrBuf.Bytes()))
				log.Println(string(stdoutBuf.Bytes()))
			}
			log.Println("Exit from waiting command execution")

		} else {
			log.Println("ERROR cmd.Start() failed with", err)
		}

		log.Println("Player has been terminated. Cmd was ", cmdText)
		actCh <- &omxstate.ActionDef{
			URI:    uri,
			Action: omxstate.ActTerminate,
		}

	}(cmdText, op.ChAction, uri, chstop)
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

	curURI := op.state.CurrURI
	if curURI != "" {
		log.Println("Shutting down the current player of ", curURI)
		if pp, ok := op.Providers[curURI]; ok {
			chStop := pp.GetCmdStopChannel()
			if chStop != nil {
				chStop <- struct{}{}
				pp.CloseStopChannel()
			}
			delete(op.Providers, curURI)
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
	op.execCommand(uri, cmd, prov.CreateStopChannel())

	return nil
}

func (op *OmxPlayer) listenStatus(statusCh chan *omxstate.StateOmx) {
	log.Println("Waiting for status in omxplayer")
	for {
		st := <-statusCh
		op.mutex.Lock()
		log.Println("Set OmxPlayer state ", st)
		if st.StatePlayer == omxstate.SPoff {
			k := op.state.CurrURI
			if _, ok := op.Providers[k]; ok {
				delete(op.Providers, k)
			}
			op.coDBus = nil
			op.clearTrackStatus()
		}
		op.state.CurrURI = st.CurrURI
		op.state.StateMute = st.StateMute
		op.state.StatePlayer = st.StatePlayer
		op.state.Info = st.Info
		op.mutex.Unlock()
	}
}

func (op *OmxPlayer) connectObjectDbBus() error {
	if op.coDBus != nil {
		return nil
	}
	u, err := user.Current()
	log.Println("User ", u.Username)

	// TODO figure out how multiple app instance works with dbus
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
