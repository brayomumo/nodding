package client

import (
	"log"
	"time"

	nodeRPC "github.com/brayomumo/nodding/nodes/grpc"
	"github.com/brayomumo/nodding/nodes/node"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func newClient(address string) *node.Nundu {
	return &node.Nundu{
		Address:  address,
		LastPing: time.Now().Unix(),
		IsAlive:  true,
	}
}

// This is meant to start up the node as invitee into a network
func Client(address, inviteToken, peerAddress string) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(peerAddress, opts...)
	if err != nil {
		log.Fatal("Error connecting to peer - ", err)
	}
	defer conn.Close()
	clientStub := nodeRPC.NewNodeClient(conn)

	client := newClient(address)
	client.HandleInvite(clientStub, inviteToken, peerAddress)
	if client.IsAlive{
		go client.SentHearbeats(clientStub, peerAddress)
	}
	for {
		time.Sleep(5*time.Second)
		client.HandleSyncLog(clientStub, "This is a beat message")
	}

}
