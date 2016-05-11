package rpcbs

import (
	"io"
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
