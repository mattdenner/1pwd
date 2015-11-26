// +build darwin

package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pborman/uuid"
)

const startItermCoprocess = `
tell application "iTerm"
	activate

	tell current window
		set cursess to (current session)
		tell cursess
			set pane to (split horizontally with same profile)
			tell pane
				log "" & id
				write text "1pwd search --coprocess={{CODE}}"
				select
			end tell
		end tell
	end tell
end tell
`

const stopItermCoprocess = `
tell application "iTerm"
	activate

	repeat with curwin in windows
		tell curwin
			repeat with curtab in tabs
				tell curtab
					repeat with curses in sessions
						tell curses
							set sid to ("" & id)
							if sid is "{{PANE}}" then
								close
							end if
						end tell
					end repeat
				end tell
			end repeat
		end tell
	end repeat
end tell
`

func doCoprocess(typ string) {
	switch typ {
	case "iterm":
		assert(runItermCoprocess())
	default:
		abortf("unable to run coprocess: %q", typ)
	}
}

func runItermCoprocess() error {
	code := uuid.New()

	l, err := net.Listen("unix", "/tmp/"+code+".sock")
	if err != nil {
		return err
	}
	defer l.Close()

	err = os.Chmod("/tmp/"+code+".sock", 0600)
	if err != nil {
		return err
	}

	script := strings.Replace(startItermCoprocess, "{{CODE}}", code, -1)
	paneID, err := exec.Command("osascript", "-e", script).CombinedOutput()
	if err != nil {
		return err
	}

	defer func() {
		id := strings.TrimSpace(string(paneID))
		if id != "" {
			script := strings.Replace(stopItermCoprocess, "{{PANE}}", id, -1)
			exec.Command("osascript", "-e", script).Run()
		}
	}()

	l.(*net.UnixListener).SetDeadline(time.Now().Add(5 * time.Minute))

	conn, err := l.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()

	l.Close()

	var data struct {
		Password string
	}

	err = json.NewDecoder(conn).Decode(&data)
	if err != nil {
		return err
	}

	if data.Password != "" {
		fmt.Print(data.Password)
	}

	return nil
}
