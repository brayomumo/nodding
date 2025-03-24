package server

import (
	"log"
	"net"

	nodeRPC "github.com/brayomumo/nodding/nodes/grpc"
	"github.com/brayomumo/nodding/nodes/node"
	"google.golang.org/grpc"
)



// This is meant to prep anything related to the gRPC server
func newServer(address string)*node.NunduServer{
	server:=  &node.NunduServer{
		Peers: make(map[string]node.Nundu),
		JoinToken: "somethinginsane1",
		Address: address,
	}
	log.Printf("Node Up and ready. Join Token: %s\nAddress: %s", server.Address, server.JoinToken)
	return server
}


func Server(address string){
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	var opt []grpc.ServerOption
	grpcServer := grpc.NewServer(opt...)
	nodeRPC.RegisterNodeServer(grpcServer, newServer(address))
	grpcServer.Serve(lis)
}