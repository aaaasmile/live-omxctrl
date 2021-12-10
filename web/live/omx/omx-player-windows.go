//+build windows

package omx

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/aaaasmile/live-omxctrl/web/live/omx/omxstate"
)

func (op *OmxPlayer) execCommand(appcmd, cmdParam, uri string, chstop chan struct{}) {
	log.Println("Prepare to start the player with execCommand")
	go func(appcmd string, cmdParam string, actCh chan *omxstate.ActionDef, uri string, chstop chan struct{}) {

		var args []string
		cmdstr := "cmd"
		appcmd = strings.ReplaceAll(appcmd, "'", "")
		args = []string{"/c", appcmd} // do not use /start

		ss1 := strings.ReplaceAll(cmdParam, "\"", "")
		paramsPart := strings.Split(ss1, " ")
		for _, ss := range paramsPart {
			args = append(args, ss)
		}

		log.Println("WINDOWS Submit the command in background ", cmdstr, args)
		cmd := exec.Command(cmdstr, args...)

		actCh <- &omxstate.ActionDef{
			URI:    uri,
			Action: omxstate.ActPlaying,
		}

		var stdoutBuf, stderrBuf bytes.Buffer
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

		//cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} // TODO windows

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
				// if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
				// 	log.Println("Error on killing the process ", err)
				// } // TODO windows
			case err := <-done:
				log.Println("Process finished")
				if err != nil {
					log.Println("Error on process termination =>", err)
				}
				log.Println(stderrBuf.String())
				log.Println(stdoutBuf.String())
			}
			log.Println("Exit from waiting command execution")

		} else {
			log.Println("ERROR cmd.Start() failed with", err)
		}

		log.Println("Player has been terminated. Cmd was ", appcmd, cmdParam)
		actCh <- &omxstate.ActionDef{
			URI:    uri,
			Action: omxstate.ActTerminate,
		}

	}(appcmd, cmdParam, op.ChAction, uri, chstop)
}
