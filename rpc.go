package main

import (
	//"fmt"
	"io"
	"net/rpc"
	"os"
	"os/exec"
	//"time"
	//"strconv"
)

type RWBridge struct {
	reader io.ReadCloser
	writer io.WriteCloser
}

func (self *RWBridge) Read(b []byte) (int, error) {
	return self.reader.Read(b)
}

func (self *RWBridge) Write(b []byte) (int, error) {
	return self.writer.Write(b)
}

func (self *RWBridge) Close() error {
	self.writer.Close()
	return self.reader.Close()
}

type Server int

func (self *Server) Ping(dummy int, reply *string) error {
	*reply = "Pong"
	return nil
}

func ServerRPC() {
	conn := &RWBridge{os.Stdin, os.Stdout}

	rpc.Register(new(Server))
	rpc.ServeConn(conn)
	return
}

func ClientRPC(command string, args ...string) (*rpc.Client, error) {
	cmd := exec.Command(command, args...)

	rx, _ := cmd.StdoutPipe()
	tx, _ := cmd.StdinPipe()
	conn := &RWBridge{rx, tx}

	cmd.Stderr = os.Stderr

	client := rpc.NewClient(conn)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	return client, nil
}

//client.Call("R.Ping", 11, &reply)
