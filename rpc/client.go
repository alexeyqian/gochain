package rpc

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// a thin wrapper around a tcp connection and golang's json-rpc client
type Client struct {
	conn    net.Conn
	jsonrpc *rpc.Client
}

func NewClient(port int) (*Client, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return nil, err
	}

	c := jsonrpc.NewClient(conn)
	client := &Client{
		conn:     conn,
		jsponrpc: c,
	}

	return client, nil
}

func (c Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	return c.jsonrpc.Call(serviceMethod, args, reply)
}

func (c Client) Close() {
	c.conn.Close()
}

/* sample code
c, err := rpc.NewClient(jsonrpcPort)
defer c.Close()

var reply string
if err := c.Call("RPC.GetMempool", nil, &reply); err != nil {
	return err
}

fmt.Println(reply)
*/
