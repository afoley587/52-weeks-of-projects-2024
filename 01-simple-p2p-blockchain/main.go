package main

// HELPFUL REFERENCES FOR ANYONE FOLLOWING ALONG
// https://blog.logrocket.com/build-blockchain-with-go/
// https://github.com/libp2p/go-libp2p/blob/master/examples/chat/README.md
// https://mycoralhealth.medium.com/code-a-simple-p2p-blockchain-in-go-46662601f417

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"os"
	"sync"
)

/*
 * We initialize a global mutex for the blockchain
 * and then we initialize the blockchain
 */
var mutex = &sync.Mutex{}
var mychain = NewBlockchain(3)

/*
 * We then pull out the flags passed in by the user (see the running section below).
 * We use our `makeHost` function to create a new node. If this is running without
 * any predefined nodes (no -d flag), then we use the `startPeer` function to wait for
 * and handle incoming streams. If there is a predefined node (there is a -d flag),
 * we initiate that connection and begin the reading/writing from the predefined node.
 * Finally, we just wait forever for the program to terminate.
 */
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sourcePort := flag.Int("sp", 0, "Source port number")
	dest := flag.String("d", "", "Destination multiaddr string")
	help := flag.Bool("help", false, "Display help")
	debug := flag.Bool("debug", false, "Debug generates the same node ID on every execution")

	flag.Parse()

	if *help {
		fmt.Printf("This program demonstrates a simple p2p blockchain application\n\n")
		fmt.Println("Usage: Run './simple-blockchain -sp <SOURCE_PORT>' where <SOURCE_PORT> can be any port number.")
		fmt.Println("Now run './simple-blockchain -d <MULTIADDR>' where <MULTIADDR> is multiaddress of previous listener host.")

		os.Exit(0)
	}

	// If debug is enabled, use a constant random source to generate the peer ID. Only useful for debugging,
	// off by default. Otherwise, it uses rand.Reader.
	var r io.Reader
	if *debug {
		// Use the port number as the randomness source.
		// This will always generate the same host ID on multiple executions, if the same port number is used.
		// Never do this in production code.
		r = mrand.New(mrand.NewSource(int64(*sourcePort)))
	} else {
		r = rand.Reader
	}

	h, err := makeHost(*sourcePort, r)
	if err != nil {
		log.Println(err)
		return
	}

	if *dest == "" {
		startPeer(ctx, h, handleStream)
	} else {
		rw, err := startPeerAndConnect(ctx, h, *dest)
		if err != nil {
			log.Println(err)
			return
		}

		// Create a thread to read and write data.
		go writeData(rw)
		go readData(rw)

	}

	// Wait forever
	select {}
}
