package peer

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"github.com/Boice-Technology/P2PNetworkProtocol/utils"
)

type MessageToSend struct {
	Content     string
	ReceiverAdd string
}

func (p *ChatPeer) readData(r *bufio.Reader) {
	for {
		str, e := r.ReadString('\n')
		if e != nil {
			log.Println(e)
		}
		str = strings.TrimSuffix(str, "\n")
		i := strings.Index(str, ":")
		if i == -1 {
			// the peer did not send its multiaddr; this should never happen
			continue
		}
		senderAdd := str[:i]
		msg := str[i+1:]
		if msg == utils.END {
			// client has said "bye bye"
			delete(p.peers, senderAdd)
			continue
		}
		p.writeFunc(msg, senderAdd)
	}
}

func (p *ChatPeer) handleWrites() {
	for messageToSend := range p.readFrom {
		w, ok := p.peers[messageToSend.ReceiverAdd]
		if !ok {
			err := p.ConnectTo(messageToSend.ReceiverAdd)
			if err != nil {
				log.Println(err)
				continue
			}
			w = p.peers[messageToSend.ReceiverAdd]
		}
		p.writeData(w, messageToSend.Content)
	}
}

func (p *ChatPeer) writeData(w *bufio.Writer, msg string) {
	w.WriteString(fmt.Sprintf("%s:%s\n", p.GetMultiaddr(), msg))
	w.Flush()
}
