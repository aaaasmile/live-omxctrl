//+build windows

package omx

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/aaaasmile/live-omxctrl/web/live/omx/omxstate"
)

func (op *OmxPlayer) execCommand(uri, cmdText string, chstop chan struct{}) {
	log.Println("Prepare to start the player with execCommand")
	go func(cmdText string, actCh chan *omxstate.ActionDef, uri string, chstop chan struct{}) {
		log.Println("WINDOWS Submit the command in background ", cmdText)
		var args []string
		cmdstr := "cmd"
		args = []string{"/c"} // do not use /start

		args = append(args, "C:\\Program Files\\VideoLAN\\VLC\\vlc.exe")
		args = append(args, "-I")
		args = append(args, "dummy")
		args = append(args, "--dummy-quiet")
		args = append(args, `c:\local\Music/CafeDelMar/cafedelmar_01.mp3`)
		cmd := exec.Command(cmdstr, args...) //"'C:\\Program Files\\VideoLAN\\VLC\\vlc.exe'", "-I", "dummy", "--dummy-quiet", "c:/local/Music/CafeDelMar/cafedelmar_01.mp3")
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

		log.Println("Player has been terminated. Cmd was ", cmdText)
		actCh <- &omxstate.ActionDef{
			URI:    uri,
			Action: omxstate.ActTerminate,
		}

	}(cmdText, op.ChAction, uri, chstop)
}
