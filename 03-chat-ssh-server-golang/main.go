package main

// ssh -o StrictHostKeyChecking=no -P 2222 localhost
// example taken from https://github.com/gliderlabs/ssh/blob/master/_examples/ssh-pty/pty.go
// as a placeholder
import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/afoley587/52-weeks-of-projects/03-chat-ssh-server-golang/rooms"
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

var sessions map[ssh.Session]*rooms.Room
var availableRooms []*rooms.Room

func chat(s ssh.Session) {
	term := terminal.NewTerminal(s, fmt.Sprintf("%s > ", s.User()))
	for {
		// get user input
		line, err := term.ReadLine()
		if err != nil {
			break
		}

		var commandPat = regexp.MustCompile(`^/enter.*`)
		if len(line) > 0 {
			if string(line[0]) == "/" {
				switch {
				case line == "/exit":
					return
				case line == "/list":
					for _, r := range availableRooms {
						term.Write([]byte(r.Name + "\n"))
					}
				case commandPat.MatchString(string(line)):
					log.Println("In /enter")
					toEnter := strings.Split(line, "/enter ")[1]
					for _, r := range availableRooms {
						if toEnter == r.Name {
							log.Println("Enterring Room...")
							r.Enter(s)
							sessions[s] = r
							log.Println(sessions)
						}
					}

				case line == "/help":
					term.Write([]byte("hey\n"))
				default:
					term.Write([]byte("hey\n"))
				}
				continue
			} else {
				sessions[s].SendMessage(line)
			}
		}
	}
}

func main() {
	availableRooms = []*rooms.Room{
		&rooms.Room{Name: "Room_A"},
		&rooms.Room{Name: "Room_B"},
		&rooms.Room{Name: "Room_C"},
	}
	sessions = make(map[ssh.Session]*rooms.Room)
	ssh.Handle(func(s ssh.Session) {
		chat(s)
	})

	log.Println("starting ssh server on port 2222...")
	log.Fatal(ssh.ListenAndServe(":2222", nil))
}
