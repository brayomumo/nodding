package node

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/brayomumo/nodding/nodes/grpc"
)

// Type to include node metadata
type Nundu struct {
	Address string
	LastPing int64
	IsAlive bool
}

type NunduServer struct {
	pb.UnimplementedNodeServer
	Peers map[string]Nundu

	// some sort of secret key used to invite other tokens
	JoinToken string
	Address string
}

// Logic for handling other peers join node network
func (n *NunduServer) Invite(_ context.Context, invite *pb.InviteGreeting) (*pb.InviteStatus, error) {
	
	log.Printf("new Node requesting to join the network: %s...", invite.Address)
	// default return status
	greetingStatus := pb.InviteStatus{
		Greeting:  invite,
		Status:    "DENIED",
		Timestamp: time.Now().Unix(),
	}
	timestamp := time.Now().Unix()
	if invite.InviteToken == n.JoinToken {
		n.Peers[invite.Address] = Nundu{
			Address: invite.Address,
			LastPing: timestamp, IsAlive: true,}
		greetingStatus = pb.InviteStatus{
			Greeting:  invite,
			Status:    "ACCEPTED",
			Timestamp: timestamp,
		}
	}
	return &greetingStatus, nil
}

// Logic to handle heartbeats from the other nodes
func (n *NunduServer) Heatbeat(stream pb.Node_HeatbeatServer) error {
	for {
		beat, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		peer := beat.Address
		peerNode, found := n.Peers[peer]
		if !found{
			// ignore heartbeats from unknown nodes
			continue
		}
		// Update node status
		peerNode.LastPing = beat.Timestamp
		peerNode.IsAlive = true
		n.Peers[peer] = peerNode
		log.Printf("heartBeat from Node: %s - %d", peer, beat.Timestamp)
		stream.Send(&pb.Beat{
			Address: n.Address,
			Timestamp: time.Now().Unix(),
		})
	}
}

// Logic for handling syncLog requests from other nodes
// Currently accepts all logs
func(n *NunduServer) SyncLog(_  context.Context, logg *pb.Log)(*pb.LogSyncStatus, error){
	log.Printf("received Log - %s - %d", logg.Message, logg.Timestamp)
	return &pb.LogSyncStatus{
		Log: logg,
		Status: "ACCEPTED",
		Timestamp: time.Now().Unix(),
	}, nil
}
