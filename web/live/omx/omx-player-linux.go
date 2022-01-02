//+build !windows

package omx

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/aaaasmile/live-omxctrl/web/live/omx/omxstate"
)

func (op *OmxPlayer) execCommand(appcmd, cmdParam, uri string, moreargs []string, chstop chan struct{}) {
	log.Println("Prepare to start the player with execCommand")
	go func(appcmd, cmdParam string, actCh chan *omxstate.ActionDef, uri string, moreargs []string, chstop chan struct{}) {
		strmore := strings.Join(moreargs, " ")
		cmdText := fmt.Sprintf("%s %s %s", appcmd, cmdParam, strmore)
		log.Println("[SUBMIT] the command in background: ", cmdText)
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
				log.Println("Process finished, error ans stdout buffers:")
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

		log.Println("Player has been TERMINATED. Cmd was ", cmdText)
		actCh <- &omxstate.ActionDef{
			URI:    uri,
			Action: omxstate.ActTerminate,
		}

	}(appcmd, cmdParam, op.ChAction, uri, moreargs, chstop)
}
