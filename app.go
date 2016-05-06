package main
import (
	"fmt"
	"time"
)


func App() int {

	for !SigINT {
		fmt.Println("hw")
		time.Sleep(time.Second * 5)
	}

	return 0
}