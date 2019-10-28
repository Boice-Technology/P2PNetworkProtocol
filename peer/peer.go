package peer

import (
	"context"
	"crypto/rand"
	"fmt"
	"bufio"

	"github.com/libp2p/go-libp2p"
	libHost "github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/libp2p/go-libp2p-core/network"
)

type OnMessageReceived func(string, string)

type ChatPeer struct {
	ip         string
	sourcePort int
	readFrom   chan MessageToSend
	writeFunc  OnMessageReceived
	p2pHost    libHost.Host
	peers      map[string]*bufio.Writer
}

func NewChatPeer(ip string, sourcePort int, readFrom chan MessageToSend, writeFunc OnMessageReceived) (*ChatPeer, error) {
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		return nil, err
	}

	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", ip, sourcePort))
	
	host, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		return nil, err
	}
	
	chatPeer := ChatPeer{ip, sourcePort, readFrom, writeFunc, host, make(map[string]*bufio.Writer)}

	// Set a function as stream handler.
	// This function is called when a peer connects, and starts a stream with this protocol.
	host.SetStreamHandler("/chat/1.0.0", func (s network.Stream) {
		// Create a buffer stream for non blocking read
		go (&chatPeer).readData(bufio.NewReader(s))
	})

	go func() {
		(&chatPeer).handleWrites()
		host.RemoveStreamHandler("/chat/1.0.0")
		host.Close()
	}()

	return &chatPeer, nil
}

func (p *ChatPeer) ConnectTo(dest string) error {
	maddr, err := multiaddr.NewMultiaddr(dest)
	if err != nil {
		return err
	}
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return err
	}

	// Add the destination's peer multiaddress in the peerstore.
	// This will be used during connection and stream creation by libp2p.
	p.p2pHost.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	// Start a stream with the destination.
	// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
	ctx, _ := context.WithCancel(context.Background())
	// // cancel() TODO: Decide if we want to cancel
	s, err := p.p2pHost.NewStream(ctx, info.ID, "/chat/1.0.0")
	if err != nil {
		return err
	}

	w := bufio.NewWriter(s)
	p.peers[dest] = w

	return nil
}

func (p *ChatPeer) GetMultiaddr() string {
	return fmt.Sprintf("/ip4/%s/tcp/%v/p2p/%s", p.ip, p.sourcePort, p.p2pHost.ID().Pretty())
}

func (p *ChatPeer) GetHostId() peer.ID {
	return p.p2pHost.ID()
}
