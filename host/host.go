package host

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	libHost "github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"

	"github.com/multiformats/go-multiaddr"
)

type WriteFunc func(string)

type ChatHost struct {
	ip         string
	sourcePort int
	readFrom   chan string
	writeFunc  WriteFunc
	p2pHost    libHost.Host
}

func Server(ip string, sourcePort int, readFrom chan string, writeFunc WriteFunc) *ChatHost {
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		panic(err)
	}

	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", ip, sourcePort))
	
	host, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		panic(err)
	}
	chatHost := ChatHost{ip, sourcePort, readFrom, writeFunc, host}

	// Set a function as stream handler.
	// This function is called when a peer connects, and starts a stream with this protocol.
	host.SetStreamHandler("/chat/1.0.0", chatHost.handleStream)

	return &chatHost
}

func Client(ip string, sourcePort int, dest string, readFrom chan string, writeFunc WriteFunc) *ChatHost {
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		panic(err)
	}

	sourceMultiaddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", ip, sourcePort))

	host, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrs(sourceMultiaddr),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		panic(err)
	}
	chatHost := ChatHost{ip, sourcePort, readFrom, writeFunc, host}
	
	// fmt.Println("This node's multiaddresses:")
	// for _, la := range host.Addrs() {
	// 	fmt.Printf(" - %v\n", la)
	// }
	// fmt.Println()

	maddr, err := multiaddr.NewMultiaddr(dest)
	if err != nil {
		log.Fatalln(err)
	}
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Fatalln(err)
	}

	// Add the destination's peer multiaddress in the peerstore.
	// This will be used during connection and stream creation by libp2p.
	host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	// Start a stream with the destination.
	// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
	s, err := host.NewStream(context.Background(), info.ID, "/chat/1.0.0")
	if err != nil {
		panic(err)
	}

	chatHost.handleStream(s)
	return &chatHost
}

func (h *ChatHost) GetMultiAddr() string {
	return fmt.Sprintf("/ip4/%s/tcp/%v/p2p/%s", h.ip, h.sourcePort, h.p2pHost.ID().Pretty())
}
