package rooms

import (
	"io"

	"github.com/gliderlabs/ssh"
)

type Room struct {
	Name     string
	History  []string
	Sessions []ssh.Session
}

func send(sess ssh.Session, message string) {
	message = message + "\n"
	io.WriteString(sess, message)
}

func (r *Room) SendMessage(message string) {
	r.History = append(r.History, message)
	for _, s := range r.Sessions {
		send(s, message)
	}
}

func (r *Room) Enter(sess ssh.Session) {
	r.Sessions = append(r.Sessions, sess)
	send(sess, "Welcome To My Room")
	for _, m := range r.History {
		send(sess, m)
	}
}
