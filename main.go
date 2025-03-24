package main

import (
	"flag"
	"fmt"

	"github.com/brayomumo/nodding/nodes/client"
	"github.com/brayomumo/nodding/nodes/server"
)

var (
	port       = flag.Int("port", 50051, "The server port")
	peerAddress = flag.String("peer_address", "", "Address of the peer inviting this node")
	joinToken = flag.String("invote_token", "", "Invite Token from the peer")

)


func main(){
	flag.Parse()
	address := fmt.Sprintf("localhost:%d", *port)
	if *peerAddress == "" || *joinToken == ""{
		server.Server(address)
	}
	client.Client(address, *joinToken, *peerAddress)
}