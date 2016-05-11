package rpcbs

import (
	"net/rpc"
	"os"
)


func Serve() {
	conn := &RWBridge{os.Stdin, os.Stdout}

	rpc.Register(new(Server))
	rpc.ServeConn(conn)
	return
}

type Server int

func (self *Server) Ping(dummy int, reply *string) error {
	*reply = "Pong"
	return nil
}

