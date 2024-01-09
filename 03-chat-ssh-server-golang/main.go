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

/*
First, we will want to keep track of a few things:
1. The sessions we have and which rooms they are in. This will help us make sure
that the user is sending messages to the right room and will also make sure that
a user leaves their current room before hopping to another room.
2. The available rooms that we can offer to users. We will be able to show users
which rooms are available, but also helps us keep track of who is where.

We also define a few supported commands:
* /enter <room> will allow a user to enter a room
* /help will show users a help screen
* /exit will exit the current session
* /list will show the chat rooms to the user
*/
var (
	sessions       map[ssh.Session]*rooms.Room
	availableRooms []*rooms.Room
	enterCmd       = regexp.MustCompile(`^/enter.*`)
	helpCmd        = regexp.MustCompile(`^/help.*`)
	exitCmd        = regexp.MustCompile(`^/exit.*`)
	listCmd        = regexp.MustCompile(`^/list.*`)
)

/*
We will use this help message when a user inputs an
unknown /command or if they type /help.
*/
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

/*
This is a function to filter a struct on some conditions
and then return the matching structs. For exmaple, we
can filter all of the rooms in our availableRooms by room
name, which is what we do below.
*/
func filter[T any](s []T, cond func(t T) bool) []T {
	res := []T{}
	for _, v := range s {
		if cond(v) {
			res = append(res, v)
		}
	}
	return res
}

/*
We will call listRooms when a user runs the /list
command. It will output each of the room's names, separated
by newlines.
*/
func listRooms() string {
	var sb strings.Builder
	for _, r := range availableRooms {
		sb.WriteString(r.Name + "\n")
	}
	return sb.String()
}

/*
chat is the main entrypoint to our chat server. First, we open a new
terminal object for a user's session where the prompt is their username
followed by a >.

We then read a line from the user with the ReadLine() function. If
the line starts with a "/", we assume its a slash command and will respond
with one of the functions we discussed above. The exception is the /enter
command. When a user tries to enter a room, we will get the room object
from our availableRooms and call the Enter function for that room which was
discussed previously (in rooms.go). To recap, that function adds a user to its
list of registered users and then shows them the entire history of the chatroom.
If they are already in a room, we first exit the old room before enterring the
new one.

If the line doesn't start with a "/", we assume that the user is trying to send
a message. If they are in a room, we simply send a message. If they are not
in a room, we just show the help command.
*/
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
				if sessions[s] != nil {
					sessions[s].SendMessage(s.User(), line)
				} else {
					term.Write([]byte((helpMsg())))
				}

			}
		}
	}
}

/*
Finally, we just wrap it all together by in our main().
First, we create three rooms. Next, we tell the ssh library
to handle new connections with the chat() function. Finally,
we use ListenAndServe to start the ssh server on port 2222.
*/
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
