//+build windows

package omx

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/aaaasmile/live-omxctrl/web/live/omx/omxstate"
	"github.com/shirou/gopsutil/v3/process"
)

func (op *OmxPlayer) execCommand(appcmd, cmdParam, uri string, moreargs []string, chstop chan struct{}) {
	log.Println("Prepare to start the player with execCommand")
	go func(appcmd string, cmdParam string, actCh chan *omxstate.ActionDef, uri string, moreargs []string, chstop chan struct{}) {

		cmd := exec.Command("cmd")
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		// Remember /c only preserve 2 quotes, so us /s
		//cmd.SysProcAttr.CmdLine = `cmd.exe /c "C:\Program Files\VideoLAN\VLC\vlc.exe" -I dummy --dummy-quiet D:\Music\ipod\883\883_casa_albergo.mp3`
		// Note that we need double quotes 2 time. The first for the program and the second for the URI
		// So we need /s in order to preserve quotes, otherwise are removed when more then 2 are written.
		// Test it with a command line, not a powershell.
		// This finally works:
		//cmd.SysProcAttr.CmdLine = `cmd.exe /s /c ""C:\Program Files\VideoLAN\VLC\vlc.exe" -I dummy --dummy-quiet "D:\Music\ipod\Bruce Springsteen - Greatest Hits Essentials 3CD [Bubanee]\CD3\06 - Missing.mp3""`
		var ccss string
		if len(moreargs) == 1 {
			ccss = fmt.Sprintf("cmd.exe /s /c \"\"%s\" %s %s \"", appcmd, cmdParam, moreargs[0])
		} else {
			log.Println("Ignore following more arguments: ", moreargs)
			ccss = fmt.Sprintf("cmd.exe /s /c \"\"%s\" %s  \"", appcmd, cmdParam)
		}
		cmd.SysProcAttr.CmdLine = ccss // avoid the mess with automatic quote in cmd.exe and exec.Command(cmd,args...)
		log.Println("WINDOWS Submit the command in background ", ccss)

		actCh <- &omxstate.ActionDef{
			URI:    uri,
			Action: omxstate.ActPlaying,
		}

		var stdoutBuf, stderrBuf bytes.Buffer
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

		if err := cmd.Start(); err == nil {
			log.Println("PID started ", cmd.Process.Pid)
			done := make(chan error, 1)
			go func() {
				done <- cmd.Wait()
				log.Println("Wait ist terminated")
			}()

			select {
			case <-chstop:
				log.Println("Received stop signal, kill parent and child processes ", cmd.Process.Pid)
				p, err := process.NewProcess(int32(cmd.Process.Pid))
				if err != nil {
					log.Println("Error on on getting process ", err)
				} else {
					if ps, err := p.Children(); err == nil {
						for _, v := range ps {
							nnc, _ := v.Exe()
							if err := v.Kill(); err != nil {
								log.Println("Error on killing the child process ", err)
							} else {
								log.Println("Killed child ", nnc)
							}
						}
					}
					if err := p.Kill(); err != nil {
						log.Println("Error on killing the main process ", err)
					}
					p.Terminate()
					nn, _ := p.Exe()
					log.Println("Killed main ", nn)
				}

			case err := <-done:
				log.Println("Process finished")
				if err != nil {
					log.Println("Error on process termination =>", err)
				}
				log.Println(stderrBuf.String())
				log.Println(stdoutBuf.String())
			}
			log.Println("Exit from execCommand")

		} else {
			log.Println("ERROR cmd.Start() failed with", err)
		}

		log.Println("Player has been terminated. Cmd was ", appcmd, cmdParam, moreargs)
		actCh <- &omxstate.ActionDef{
			URI:    uri,
			Action: omxstate.ActTerminate,
		}

	}(appcmd, cmdParam, op.ChAction, uri, moreargs, chstop)
}
