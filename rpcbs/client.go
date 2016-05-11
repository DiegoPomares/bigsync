package rpcbs

import (
	"fmt"
	"net/rpc"
	"os"
	"os/exec"
	"time"
)


func Client(args ...string) (*client, error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr

	rx, _ := cmd.StdoutPipe()
	tx, _ := cmd.StdinPipe()
	bridge := &RWBridge{rx, tx}

	rpc_client := rpc.NewClient(bridge)

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	return &client{rpc_client}, nil
}

type client struct {
	rpc_client *rpc.Client
}

func (self *client) Ping() string {
	var reply string

	init_time := time.Now()

	self.rpc_client.Call("Server.Ping", 0, &reply)
	return fmt.Sprintf("%s %v", reply, time.Now().Sub(init_time))
}

func (self *client) Ping2() (string, error) {
	var reply string

	err := self.rpc_client.Call("Server.Ping", 0, &reply)
	if err != nil {
		return "", err
	}

	return reply, nil
}
