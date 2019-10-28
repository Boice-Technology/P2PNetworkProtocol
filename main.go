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
	peer1, err := peer.NewChatPeer("127.0.0.1", 5000, peer1MessageSender, func(msg string, senderId string) {
		fmt.Println("peer1 Received:", msg, "From:", senderId)

		// Handle the received message
		if msg == "Hi peer" {
			peer1MessageSender <- peer.MessageToSend{"Bye peer", senderId}
		} else {
			peer1MessageSender <- peer.MessageToSend{utils.END, senderId}
		}
	})
	if err != nil {
		log.Fatalln(err)
	}

	peer2MessageSender := make(chan peer.MessageToSend, 10)
	peer2, err := peer.NewChatPeer("127.0.0.1", 5001, peer2MessageSender, func(msg string, senderId string) {
		fmt.Println("peer2 Received:", msg, "From:", senderId)
		
		// Handle the received message
		if msg == "Hello peer" {
			peer2MessageSender <- peer.MessageToSend{"Hi peer", senderId}
		} else {
			peer2MessageSender <- peer.MessageToSend{":)", senderId}
		}
	})
	if err != nil {
		log.Fatalln(err)
	}
	err = peer2.ConnectTo(peer1.GetMultiaddr())
	if err != nil {
		log.Fatalln(err)
	}
	
	fmt.Println("Connected")

	peer1MessageSender <- peer.MessageToSend{"Hello peer", peer2.GetHostId().Pretty()}
	
	time.Sleep(2 * time.Second)
	
	close(peer1MessageSender)
	close(peer2MessageSender)
}
