package main

import (
	"fmt"
	"time"
	"log"

	"github.com/Boice-Technology/P2PNetworkProtocol/peer"
	"github.com/Boice-Technology/P2PNetworkProtocol/utils"
)

func main() {
	peer1MessageSender := make(chan peer.MessageToSend, 10)
	peer1, err := peer.NewChatPeer("127.0.0.1", 5000, peer1MessageSender, func(msg string, senderAdd string) {
		fmt.Println("peer1 Received:", msg, "From:", senderAdd)

		// Handle the received message
		if msg == "Hello peer" {
			peer1MessageSender <- peer.MessageToSend{"Hi peer", senderAdd}
		} else {
			peer1MessageSender <- peer.MessageToSend{":)", senderAdd}
		}
	})
	if err != nil {
		log.Fatalln(err)
	}

	peer2MessageSender := make(chan peer.MessageToSend, 10)
	_, err = peer.NewChatPeer("127.0.0.1", 5001, peer2MessageSender, func(msg string, senderAdd string) {
		fmt.Println("peer2 Received:", msg, "From:", senderAdd)
		
		// Handle the received message
		if msg == "Hi peer" {
			peer2MessageSender <- peer.MessageToSend{"Bye peer", senderAdd}
		} else {
			peer2MessageSender <- peer.MessageToSend{utils.END, senderAdd}
		}
	})
	if err != nil {
		log.Fatalln(err)
	}

	peer3MessageSender := make(chan peer.MessageToSend, 10)
	_, err = peer.NewChatPeer("127.0.0.1", 5002, peer3MessageSender, func(msg string, senderAdd string) {
		fmt.Println("peer3 Received:", msg, "From:", senderAdd)
		
		// Handle the received message
		if msg == "Hi peer" {
			peer3MessageSender <- peer.MessageToSend{"Bye peer", senderAdd}
		} else {
			peer3MessageSender <- peer.MessageToSend{utils.END, senderAdd}
		}
	})
	if err != nil {
		log.Fatalln(err)
	}

	peer2MessageSender <- peer.MessageToSend{"Hello peer", peer1.GetMultiaddr()}
	peer3MessageSender <- peer.MessageToSend{"Hello peer", peer1.GetMultiaddr()}
	
	time.Sleep(2 * time.Second)
	
	close(peer1MessageSender)
	close(peer2MessageSender)
	close(peer3MessageSender)
}
