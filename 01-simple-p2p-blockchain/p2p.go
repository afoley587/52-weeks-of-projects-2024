package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/kr/pretty"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
)

/*
 * Messages will come in from a user over a CLI.
 * In our system, it will just be a JSON object
 * similar to {"location": "hawaii", "waveHeight": 4}
 * But, of course, in a production system, this might come
 * from an API or data stream or some other data source
 */
type UserMessage struct {
	Location   string
	WaveHeight int
}

/*
 * First, and probably easiest, we will need to make a Peer-to-peer
 * host. Now, P2P was a bit confusing for me, but I am extremely grateful
 * for their great docs and examples (https://github.com/libp2p/go-libp2p/tree/master/examples/chat)
 * which served as a ground-work for this P2P blockchain.
 * This host will be one of the blockchain nodes in the system and other nodes
 * will be able to latch on to it's networking details and connect to it.
 *
 * The first thing we do is create an RSA key using a random seed. P2P uses
 * the public/private key pair to keep our system's secure. We then create
 * a Multiaddr that listens on 0.0.0.0 and some customizeable port passed in by
 * the user. Finally, we return a new P2P host by calling the New() function with
 * our constructed multiaddr and our created private RSA key.
 */
func makeHost(port int, randomness io.Reader) (host.Host, error) {
	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, randomness)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	return libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
}

/*
 * Next, we need a way to handle incoming streams. These streams will
 * be instantiated when another node connects to this current node. If the
 * nodes want to be able to send/receive data, we need to handle both the
 * inbound and outbound data packets of the stream.
 *
 * Our stream is bi-directional, so we need to create a ReadWriter to handle
 * both directions. We also want to handle the reading/writing concurrently,
 * so we need to put them on separate go routines.
 *
 * As a note, most of our business logic will go into those readData and
 * writeData functions. Up until now, we've more or less been following the
 * examples on the LibP2P GitHub page.
 */
func handleStream(s network.Stream) {
	log.Println("Stream detected")

	// Create a buffer stream for non-blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
	go writeData(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}

/*
 * First, let's look at how we read data from the stream and add blocks
 * to our block chain. We read a string from the stream up the a new line character.
 * If we get an EOF error, we can assume the stream is closed and return. If we get
 * an empty string, we can just skip to the next loop iteration.
 * However, if we successfully read data from the stream, then we should try to
 * unmarshal it into a blockchain object.
 *
 * In block-chain land, the length of the chain is king. If the incoming chain is longer
 * than our chain, we should assume that the new chain is the most up-to-date
 * statement of record and we should essentially discard our chain for the incoming
 * chain.
 *
 * Finally, we just use the pretty package to dump out our chain to the screen.
 */
func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')

		// If the channel is closed or we get an EOF, return
		if err == io.EOF {
			return
		}

		if str == "" {
			continue
		}

		if str != "\n" {
			var chain Blockchain
			if err := json.Unmarshal([]byte(str), &chain); err != nil {
				log.Println(err)
				continue
			}

			mutex.Lock()
			if len(chain.Chain) > len(mychain.Chain) {
				mychain = chain
				pretty.Println(mychain)
			}
			mutex.Unlock()
		}

	}
}

/*
 * We have a way to read data from the stream, but how about
 * writing data to the stream? Enter writeData. This function
 * will read data from standard input (os.Stdin). It will then
 * unmarshal it from a string to a UserMessage object (defined above).
 * It will then create a block from this message, validate it, and
 * write the updated chain to the stream.
 */
func writeData(rw *bufio.ReadWriter) {

	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		sendData = strings.Replace(sendData, "\n", "", -1)
		var userMsg UserMessage

		if err := json.Unmarshal([]byte(sendData), &userMsg); err != nil {
			log.Println(err)
			continue
		}

		mutex.Lock()
		mychain.appendBlock(userMsg.Location, userMsg.WaveHeight)
		if !mychain.isValid() {
			pretty.Println(mychain)
			log.Println("Chain isn't valid anymore! Help meee!")
			return
		}
		mutex.Unlock()

		bytes, err := json.Marshal(mychain)
		log.Println(string(bytes))
		if err != nil {
			log.Println(err)
		}

		pretty.Println(mychain)

		mutex.Lock()
		rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		rw.Flush()
		mutex.Unlock()
	}

}

/*
 * startPeer will be used to start a node on it's own. It will
 * set the stream handler for the host and then print out it's connection
 * details to stdout so that other peers can connect directly to it.
 *
 * startPeerAndConnect is similar except it will initiate connections
 * with the destination multiaddr we specify on the command line. So,
 * in this case, we assume that at least one node is already running
 * and we will connect directly to it.
 */
func startPeer(ctx context.Context, h host.Host, streamHandler network.StreamHandler) {
	// Set a function as stream handler.
	// This function is called when a peer connects, and starts a stream with this protocol.
	// Only applies on the receiving side.
	h.SetStreamHandler("/chat/1.0.0", streamHandler)

	// Let's get the actual TCP port from our listen multiaddr, in case we're using 0 (default; random available port).
	var port string
	for _, la := range h.Network().ListenAddresses() {
		if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
			port = p
			break
		}
	}

	if port == "" {
		log.Println("was not able to find actual local port")
		return
	}

	log.Printf("Run './simple-blockchain -d /ip4/127.0.0.1/tcp/%v/p2p/%s' on another console.\n", port, h.ID())
	log.Println("You can replace 127.0.0.1 with public IP as well.")
	log.Println("Waiting for incoming connection")
	log.Println()
}

func startPeerAndConnect(ctx context.Context, h host.Host, destination string) (*bufio.ReadWriter, error) {
	log.Println("This node's multiaddresses:")
	for _, la := range h.Addrs() {
		log.Printf(" - %v\n", la)
	}
	log.Println()

	// Turn the destination into a multiaddr.
	maddr, err := multiaddr.NewMultiaddr(destination)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Add the destination's peer multiaddress in the peerstore.
	// This will be used during connection and stream creation by libp2p.
	h.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	// Start a stream with the destination.
	// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
	s, err := h.NewStream(context.Background(), info.ID, "/chat/1.0.0")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("Established connection to destination")

	// Create a buffered stream so that read and writes are non-blocking.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	return rw, nil
}
