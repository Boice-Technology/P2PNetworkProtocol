package utils

import (
	"github.com/libp2p/go-libp2p-core/peer"
)

func prettyIdToPeerId(peerId string) (peer.ID, error) {
	return peer.IDB58Decode(peerId)
}
