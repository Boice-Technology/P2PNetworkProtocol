package host

import (
	"bufio"
	"fmt"

	"github.com/libp2p/go-libp2p-core/network"
)

func (h *ChatHost) handleStream(s network.Stream) {
	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go h.readData(rw)
	go h.writeData(rw)
	// the stream will stay open until you close it (or the other side closes it).
}

func (h *ChatHost) readData(rw *bufio.ReadWriter) {
	for {
		str, e := rw.ReadString('\n')
		if e != nil {
			break
		}
		// if str == "" {
		// 	return
		// }
		if str != "\n" {
			h.writeFunc(str)
		}
	}
}

func (h *ChatHost) writeData(rw *bufio.ReadWriter) {
	for sendData := range h.readFrom {
		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()
	}
}
