package rpcbs

import (
	"fmt"
	"net/rpc"
	"os"
	"os/exec"
	"time"
	//"github.com/DiegoPomares/bigsync/hasher"
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


func (self *client) SetParams(dest_file string, file_size int64, block_size,
	workers int, hash_type string, force_creation bool, mode string) error {

	data := Params{dest_file, file_size, block_size, workers, hash_type, force_creation, mode}

	return self.rpc_client.Call("Server.SetParams", data, nil)
}

func (self *client) StartHashing() error {

	return self.rpc_client.Call("Server.StartHashing", 0, nil)
}
