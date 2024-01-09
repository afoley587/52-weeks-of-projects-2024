package rooms

import (
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

/*
First, let us look at the structs that our rooms package will use
to keep track of data. In our system, a chatroom will be an object
which keeps track of users in the chatroom, the messages that have
been sent, and will help facilitate the delivery of messages to all
of the users in the room.

The Room struct will have a few pieces of data associated with it.
 1. The name of the room: This will be a string identifier to show
    to users when they use the /enter command.
 2. The chat history of the room. The room will keep all of its chat
    history in memory. When a new user joins, it can then show them
    the entire history of the chatroom so they can feel caught-up.
 3. The users currently in the chatroom. This way, the system can
    notify the users of new messages and keep a record of who is
    currently logged in to the room.

As previously said, the room needs to keep track of two things:
users and messages. So, what's in a user? Well, our User struct
will be comprised of an SSH Session and the associated PTY terminal.
The SSH Session gives us underlying session information like the
username while the terminal gives us a nice way to send and
receive messages from the user. The message struct, on the other hand,
just needs to keep track of the sending user and the actual message.
*/
type Room struct {
	Name    string
	History []Message
	Users   []User
}

type User struct {
	Session  ssh.Session
	Terminal *terminal.Terminal
}

type Message struct {
	From    string
	Message string
}

/*
So, now that our room is established, what kind of things should the room
be able to do? Well, I think our room should allow users to do three things:
1. Enter the room
2. Leave the room
3. Send messages to the room

Let's first look at our Enter function. This is the function that will be
called when a user uses the /enter command after logging in. We see
that the first thing we do is make a User object out the users
session and terminal. The observant reader will notice we don't do
any checking to see if the user is already in the room. This is mainly
because we force users to leave rooms before joining new ones, and this
is enforced one layer up.

The second thing our Enter function does is add the new user to the
room's user slice. Then the room sends them a message and sends them
all of the messages in its history.

Let us turn our attention to leaving a room. Once a user has decided
that they want to go to another chatroom, they will call the Leave
function to remove them from this current room. All our leave function
has to do is remove them from the user slice. Note that we remove
users by username so, if two people were to be logged into the same
account, some funky behavior might happen. As of now, the system
doesn't enforce that but, it would be a good idea to!

Finally, a user has to be able to send messages to other users. They
do that by typing in their terminal. First, we create a message
object with their username as the From attribute and the string
they typed as the Message attribute. The message then get's appended
to the rooms history and sent out to all of the other users.
*/

func (r *Room) Enter(sess ssh.Session, term *terminal.Terminal) {
	u := User{Session: sess, Terminal: term}
	r.Users = append(r.Users, u)
	entryMsg := Message{From: r.Name, Message: "Welcome to my room!"}
	send(u, entryMsg)
	for _, m := range r.History {
		send(u, m)
	}
}

func (r *Room) Leave(sess ssh.Session) {
	r.Users = removeByUsername(r.Users, sess.User())
}

func (r *Room) SendMessage(from, message string) {

	messageObj := Message{From: from, Message: message}
	r.History = append(r.History, messageObj)
	for _, u := range r.Users {
		if (u.Session.User()) != from {
			send(u, messageObj)
		}
	}
}

func removeByUsername(s []User, n string) []User {
	var index int
	for i, u := range s {
		if u.Session.User() == n {
			index = i
			break
		}
	}
	return append(s[:index], s[index+1:]...)
}

func send(u User, m Message) {
	raw := m.From + "> " + m.Message + "\n"
	u.Terminal.Write([]byte(raw))
}
