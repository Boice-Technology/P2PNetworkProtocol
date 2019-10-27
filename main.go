package main

import (
	"fmt"
	"time"

	"github.com/Boice-Technology/P2PNetworkProtocol/host"
)

func main() {
	serverMessageSender := make(chan string, 10)
	server := host.Server("0.0.0.0", 5000, serverMessageSender, func(msg string) {
		fmt.Print("Server Received: ", msg)
	})

	clientMessageSender := make(chan string, 10)

	host.Client("0.0.0.0", 5001, server.GetMultiAddr(), clientMessageSender, func(msg string) {
		fmt.Print("Client Received: ", msg)
		
		// Handle the received message
		if msg == "Hello client\n" {
			clientMessageSender <- "Hi server"
		} else {
			clientMessageSender <- ":)"
		}
	})
	
	fmt.Println("Connected")

	serverMessageSender <- "Hello client"
	serverMessageSender <- "Bye client"
	close(serverMessageSender)
	
	time.Sleep(2 * time.Second)
	
	close(clientMessageSender)
}