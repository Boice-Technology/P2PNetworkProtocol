package peer

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/libp2p/go-libp2p-core/network"

	"github.com/Boice-Technology/P2PNetworkProtocol/utils"
)

type MessageToSend struct {
	Content    string
	ReceiverId string
}

func (p *ChatPeer) handleStream(s network.Stream) {
	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go p.readData(rw)
	go p.writeData(rw)
	// the stream will stay open until you close it (or the other side closes it).
}

func (p *ChatPeer) readData(rw *bufio.ReadWriter) {
	for {
		str, e := rw.ReadString('\n')
		if e != nil {
			break
		}
		str = strings.TrimSuffix(str, "\n")
		i := strings.Index(str, ":")
		if i == -1 {
			// the peer did not send its id; this should not happen
			continue
		}
		senderId := str[:i]
		msg := str[i+1:]
		if msg == utils.END {
			// disconnect from this peer
			continue
		}
		p.writeFunc(msg, senderId)
	}
}

// If it is a client that is writing, the stream is already created
// If it is a server, it writes only after it has received the stream, so it is already created
func (p *ChatPeer) writeData(rw *bufio.ReadWriter) {
	for messageToSend := range p.readFrom {
		rw.WriteString(fmt.Sprintf("%s:%s\n", p.GetHostId().Pretty(), messageToSend.Content))
		rw.Flush()
	}
}
