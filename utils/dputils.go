package utils

import (
	"fmt"
	"os"
	"time"
)


func Stderrln(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}

func Stderrf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}


func Iserror(err error, args ...interface{}) bool {
	if err != nil {
		Stderrf("%s. ", err)
		Stderrln(args...)
	
		return true
	}

	return false
}

func Sleep(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}