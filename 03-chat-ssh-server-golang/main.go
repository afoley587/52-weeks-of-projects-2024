package main

// ssh -o "StrictHostKeyChecking=no" -p 2222 alex@127.0.0.1

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/afoley587/52-weeks-of-projects/03-chat-ssh-server-golang/rooms"
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	sessions       map[ssh.Session]*rooms.Room
	availableRooms []*rooms.Room
	enterCmd       = regexp.MustCompile(`^/enter.*`)
	helpCmd        = regexp.MustCompile(`^/help.*`)
	exitCmd        = regexp.MustCompile(`^/exit.*`)
	listCmd        = regexp.MustCompile(`^/list.*`)
)

func helpMsg() string {
	return `
Hello and welcome to the chat server! Please use
one of the following commands:
	1. /list: To list available rooms
	2. /enter <room>: To enter a room
	3. /exit: To leave the server
	4. /help: To display this message
`
}

func filter[T any](s []T, cond func(t T) bool) []T {
	res := []T{}
	for _, v := range s {
		if cond(v) {
			res = append(res, v)
		}
	}
	return res
}

func listRooms() string {
	var sb strings.Builder
	for _, r := range availableRooms {
		sb.WriteString(r.Name + "\n")
	}
	return sb.String()
}

func chat(s ssh.Session) {
	term := terminal.NewTerminal(s, fmt.Sprintf("%s > ", s.User()))
	for {
		line, err := term.ReadLine()
		if err != nil {
			break
		}

		if len(line) > 0 {
			if string(line[0]) == "/" {
				switch {
				case exitCmd.MatchString(string(line)):
					return
				case listCmd.MatchString(string(line)):
					term.Write([]byte(listRooms()))
				case enterCmd.MatchString(string(line)):
					toEnter := strings.Split(line, " ")[1]
					matching := filter(availableRooms, func(r *rooms.Room) bool {
						return toEnter == r.Name
					})
					if len(matching) == 0 {
						term.Write([]byte("Invalid Room!\n"))
					} else {
						if sessions[s] != nil {
							sessions[s].Leave(s)
						}
						r := matching[0]
						r.Enter(s, term)
						sessions[s] = r
					}
				case helpCmd.MatchString(string(line)):
					term.Write([]byte(helpMsg()))
				default:
					term.Write([]byte((helpMsg())))
				}
			} else {
				sessions[s].SendMessage(s.User(), line)
			}
		}
	}
}

func main() {
	availableRooms = []*rooms.Room{
		&rooms.Room{Name: "a"},
		&rooms.Room{Name: "b"},
		&rooms.Room{Name: "c"},
	}
	sessions = make(map[ssh.Session]*rooms.Room)
	ssh.Handle(func(s ssh.Session) {
		chat(s)
	})

	log.Println("starting ssh server on port 2222...")
	log.Fatal(ssh.ListenAndServe(":2222", nil))
}
