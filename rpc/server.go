package rpc

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// RPC server is responsible for reading data from a connection.
// parsing it, handling it and sending a replay.
// Similar to HTTP server, but we don't need to define paths and there are no methods

// a thin wrapper around the net/rpc Server
type Server struct {
	port int
	rpc  *rpc.Server
}

//rpcsvr is the actual RPC server.
// In Golang, it’s an abstraction that doesn’t depend on transport layer:
// we can use TCP or HTTP, RPC server doesn’t need to know which we choose.
func NewServer(port int, node Node) (*Server, error) {
	rpcsvr := rpc.NewServer()

	handlers := RPC{node: node}
	if err := rpcsvr.Register(handlers); err != nil {
		return nil, err
	}

	s := Server{
		port: port,
		rpc:  rpcsvr,
	}

	return &s, nil
}

func (s Server) Run() {
	l, err := net.Listen("tcp", fmt.Sprintf("%d", s.port))
	for {
		conn, err := l.Accept()
		go s.rpc.ServerCodec(jsonrpc.NewServerCodec(conn))
	}
}
