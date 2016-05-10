package main

import (
	"fmt"
	"os"
	//"time"
	//"io"
	//"strconv"
	"strings"
	"github.com/DiegoPomares/bigsync/hasher"
	u "github.com/DiegoPomares/bigsync/utils"
)

//const REMOTE_COMMAND = "bigsync --server-mode-enable dummy"
const REMOTE_COMMAND = "go run /home/local/ANT/dieamare/dev/go/src/github.com/DiegoPomares/bigsync/*.go --server-mode-enable dummy"

var sigINT = false
func AppSignal(sig os.Signal) {
	if !sigINT && sig == os.Interrupt {
		u.Stderrln("[ Interrupt received, cleaning up ... ]")
		sigINT = true
	}

}

func App() error {

	if !Options.ServerMode {

		//TODO open file here
		fh, err := hasher.New(Options.SourceFile, "r", Options.BlockSize, Options.HashType, Options.Workers)
		if u.Iserror(err) {
			return err
		}

		if Verbose {
			u.Stderrln("File size", fh.FileSize, "| Blocks", fh.NumBlocks)
		}

		// Just print hashes of local file ------------------------------------
		if Options.DestFile == "" {

			// Start hashing
			fh.Start()

			// Print header
			fmt.Printf("{\n")
			fmt.Printf("  \"file\": \"%s\",\n", fh.FilePath)
			fmt.Printf("  \"block_size\": %d,\n", fh.BlockSize)
			fmt.Printf("  \"hash_type\": \"%s\",\n", fh.HashType)
			fmt.Printf("  \"blocks\": [\n")

			var last_read, last_block int
			//for result := range fh.Hashes {
			for result := range fh.HashRange() {

				if result.Index != 0 {
					fmt.Printf(",\n")
				}
				fmt.Printf("    { \"block\": %d, \"hash\": \"%x\" }", result.Index, result.Hash)
				last_block = result.Index
				last_read = result.Size
			}

			// Close jobs channel, wait for printer, then print footer
			fmt.Printf("\n  ],\n")
			fmt.Printf("  \"last_block\": %d,\n", last_block)
			fmt.Printf("  \"last_block_size\": %d,\n", last_read)
			fmt.Printf("  \"last_block_diff\": %d\n", Options.BlockSize-last_read)
			fmt.Printf("}\n")

			// --------------------------------------------------------------------
		} else {

			var command []string

			// Custom shell mode
			if Options.CustomSh != "" {
				command = append(command, strings.Split(Options.CustomSh, " ")...)

			// Local mode
			} else if Options.RemoteHost == "" {
				command = []string{"sh", "-c", REMOTE_COMMAND}

			// SSH mode
			} else {
				command = []string{"ssh"}
				if Options.ExtraSsh != "" {
					command = append(command, strings.Split(Options.ExtraSsh, " ")...)
				}
				command = append(command, Options.RemoteHost, REMOTE_COMMAND)

			}


			// Open connection to RPC server
			client, err := ClientRPC(command...)
			if u.Iserror(err) {
				return err
			}

			var reply string
			err = client.Call("Server.Ping", 0, &reply)
			if u.Iserror(err) {
				return err
			}
			fmt.Println(reply)

		}

	} else {
		ServerRPC()
	}

	return nil
}
