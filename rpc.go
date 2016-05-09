package main

import (
	//"fmt"
	"io"
	"net/rpc"
	"os"
	//"os/exec"
	//"time"
	//"strconv"
)

type RWBridge struct {
	reader io.ReadCloser
	writer io.WriteCloser
}

func (self *rwBridge) Read(b []byte) (int, error) {
	return self.reader.Read(b)
}

func (self *rwBridge) Write(b []byte) (int, error) {
	return self.writer.Write(b)
}

func (self *rwBridge) Close() error {
	self.writer.Close()
	return self.reader.Close()
}


type Server int

func (self *Server) Ping(n int, reply *string) error {
	*reply = "Pong"
	return nil
}


func ServerRPC() {
	conn := &rwBridge{os.Stdin, os.Stdout}

	rpc.Register(new(R))
	rpc.ServeConn(conn)
	return
}

func ClientRPC() {
	comm := exec.Command("go", "run", "rpc.go", "server")
	rx, _ := comm.StdoutPipe()
	tx, _ := comm.StdinPipe()
	conn := &rwBridge{rx, tx}

	client := rpc.NewClient(conn)
	err := comm.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
}

//client.Call("R.Ping", 11, &reply)