package main

import (
	//"fmt"
	"os"
	"os/signal"
	"runtime"
	u "github.com/DiegoPomares/bigsync/utils"
)

const (
	OK = 0
	GEN_ERROR = 1
	PARSE_ERROR = 2
	IO_ERROR = 3
	SIGINT = 130
)

var signals = make(chan os.Signal)
func process_signals() {
	for sig := range signals {
		AppSignal(sig)
	}
}

func init() {
	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())

}

func main() {

	// Process command line arguments
	err := Process_opts()
	if u.Iserror(err) {
		os.Exit(PARSE_ERROR)
	}

	// Print Version and options if verbose
	if Verbose {
		u.Stderrf("%s %s\nOptions:\n%+v\n", NAME, VERSION, Options)
	}

	// Signal handlers
	go process_signals()
	signal.Notify(signals, os.Interrupt)  // SIGINT (Ctrl+C)
	//signal.Notify(signals, syscall.SIGUSR1)  // Custom user signal 1
	//signal.Notify(signals, syscall.SIGUSR2)  // Custom user signal 2
	signal.Reset() //DEBUG


	// Run app
	status := App()


	// Exit status handling
	if _, ok := status.(*os.PathError); ok {
		os.Exit(IO_ERROR)
	}

	if status != nil {
		os.Exit(GEN_ERROR)
	}


	os.Exit(OK)
}
