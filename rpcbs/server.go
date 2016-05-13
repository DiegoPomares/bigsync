package rpcbs

import (
	"net/rpc"
	"os"
	"github.com/DiegoPomares/bigsync/hasher"
	u "github.com/DiegoPomares/bigsync/utils"
)


func Serve() {
	conn := &RWBridge{os.Stdin, os.Stdout}

	rpc.Register(new(Server))
	rpc.ServeConn(conn)
	return
}

type Server struct {
	options Params
	fhasher *hasher.Hasher
}

func (self *Server) Ping(dummy int, reply *string) error {
	*reply = "Pong"
	return nil
}

func (self *Server) SetParams(params Params, _ *int) error {
	self.options = params

	u.Stderrln(params)

	return nil
}

func (self *Server) StartHashing(_ int, _ *int) error {
	var err error

	// Open file
	mode := "r"
	if self.options.Mode == "sync" {
		mode = "rw"
	}

	self.fhasher, err = hasher.New(self.options.DestFile, mode, self.options.BlockSize,
		self.options.HashType, self.options.Workers)
	if err != nil {
		return err
	}




	return nil
}