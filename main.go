package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	u "github.com/DiegoPomares/bigsync/utils"
)

var Signals = make(chan os.Signal)
var SigINT = false

func process_signals() {
	for sig := range Signals {
		
		switch sig {
		case os.Interrupt:
			if !SigINT {
				fmt.Fprintf(os.Stderr, "[ Interrupt received, cleaning up ... ]\n")
			}
			SigINT = true
		}
	}
}

func init() {
	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())

}

func main() {

	// Process command line arguments
	err := Process_opts()
	if u.Perror(err) {
		os.Exit(1)
	}

	if Options.Verbose {
		fmt.Fprintf(os.Stderr, "%s %s\nOptions:\n%+v\n", NAME, VERSION, Options)
	}

	// Signal handlers
	signal.Notify(Signals, os.Interrupt)  // SIGINT (Ctrl+C)
	go process_signals()

	// Run app
	os.Exit(App())

}
