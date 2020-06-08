package ginger

import (
	"google.golang.org/grpc"
	"log"
	"net"
)

type GRPCContainer struct {
	server     *grpc.Server
	listerConn net.Listener
}

func (g *GRPCContainer) New(addr string) *GRPCContainer {
	var err error
	g.listerConn, err = net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	g.server = grpc.NewServer()
	//pb.RegisterGreeterServer(s, &server{})
	log.Println("rpc服务已经开启")
	return g
}

func (g *GRPCContainer) Register(server *interface{}) {
	//g.server.RegisterService(server)
	g.server.Serve(g.listerConn)
}
