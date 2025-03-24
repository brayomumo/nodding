package node

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/brayomumo/nodding/nodes/grpc"
)

func (n *Nundu) HandleInvite(client grpc.NodeClient, joinToken, peerAddress string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	status, err := client.Invite(ctx, &grpc.InviteGreeting{
		Address:     n.Address,
		InviteToken: joinToken,
		Timestamp:   time.Now().Unix(),
	})
	if err != nil {
		log.Fatal(err)
	}
	if status.Status != "ACCEPTED" {
		log.Fatal("Join Request Denied!")
	}
	n.IsAlive = true
}

func (n *Nundu) SentHearbeats(client grpc.NodeClient, peer string) {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()
	stream, err := client.Heatbeat(ctx)
	if err != nil {
		log.Fatal("Heart is not beating anymore", err)
	}
	waitChannel := make(chan struct{})
	// go routine to handle incoming heartbeats
	go func() {
		for {
			beat, err := stream.Recv()
			// loosing nodes is not catastrophic
			if err == io.EOF {
				log.Printf("Lost Node - %s", beat.Address)
				close(waitChannel)
			} else if err != nil {
				log.Fatal("client.Heartbeat failed - ", err)
			}
			log.Printf("Heartbeat from peer - %s", beat.Address)
		}
	}()

	for {
		time.Sleep(5*time.Second)
		err := stream.Send(&grpc.Beat{
			Address:   n.Address,
			Timestamp: time.Now().Unix(),
		})
		if err != nil {
			stream.CloseSend()
			<-waitChannel
			log.Fatal("Cannot send heartbeats to peers - ", err)
		}
	}
	
}


func (n *Nundu) HandleSyncLog( client grpc.NodeClient, msg string){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	status, err := client.SyncLog(ctx, &grpc.Log{
		Message: msg,
		Timestamp: time.Now().Unix(),
	})
	if err != nil{
		log.Fatal("Cannot sync logs")
	}
	log.Printf("Synced logs: Status: %s", status.Status)
}