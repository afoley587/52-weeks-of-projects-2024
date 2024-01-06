package main

// example taken from https://github.com/gliderlabs/ssh/blob/master/_examples/ssh-pty/pty.go
// as a placeholder
import (
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
	"unsafe"

	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

var sessions []io.Writer

func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

func broadcast(m string) {
	for _, s := range sessions {
		io.WriteString(s, m)
	}
}

func welcome(s ssh.Session) {
	broadcast(fmt.Sprintln("Welcome, ", s.User()))
}

func chat(s ssh.Session) {
	term := terminal.NewTerminal(s, fmt.Sprintf("%s > ", s.User()))
	for {
		// get user input
		line, err := term.ReadLine()
		if err != nil {
			break
		}

		if len(line) > 0 {
			if string(line[0]) == "/" {
				switch line {
				case "/exit":
					return
				case "/help":
					term.Write([]byte("hey"))
				default:
					term.Write([]byte("hey"))
				}
				continue
			}
		}
	}
}

func main() {
	ssh.Handle(func(s ssh.Session) {
		sessions = append(sessions, s)
		welcome(s)
		chat(s)
	})

	log.Println("starting ssh server on port 2222...")
	log.Fatal(ssh.ListenAndServe(":2222", nil))
}
