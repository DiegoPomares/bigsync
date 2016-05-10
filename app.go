package main

import (
	"fmt"
	"os"
	"time"
	"io"
	"math"
	"strconv"
	u "github.com/DiegoPomares/bigsync/utils"
)

var sigINT = false
func AppSignal(sig os.Signal) {
	if !sigINT && sig == os.Interrupt {
		u.Stderrln("[ Interrupt received, cleaning up ... ]")
		sigINT = true
	}

}

func App() error {

	if Options.ServerMode == "" {

		// Open source file
		source_file, err := os.OpenFile(Options.SourceFile, os.O_RDONLY, 0)
		if u.Iserror(err) {
			return err
		}
		defer source_file.Close()

		file_info, err := source_file.Stat()
		if u.Iserror(err) {
			return err
		}
		file_size := file_info.Size()
		num_blocks := int(math.Ceil(float64(file_size)/float64(Options.BlockSize)))

		if Verbose {
			u.Stderrln("File size", file_size, "| Blocks", num_blocks)
		}

		// Queues
		jobs := make(chan Block)
		results := make(chan Block)

		// Just print hashes of local file ------------------------------------
		if Options.RemoteHost == "" {

			StartHashWorkers(Options.Workers, jobs, results, Options.HashAlgorithm)
			StartPrinter(results)

			// Print header
			fmt.Printf("{\n")
			fmt.Printf("  \"file\": \"%s\",\n", Options.SourceFile)
			fmt.Printf("  \"block_size\": %d,\n", Options.BlockSize)
			fmt.Printf("  \"hash_type\": \"%s\",\n", Options.HashAlgorithm)
			fmt.Printf("  \"blocks\": [\n")

			lastblock, lastread, err := read_file(source_file, Options.BlockSize, jobs)
			if u.Iserror(err, "Block number", strconv.Itoa(lastblock)) {
				return err
			}

			// Close jobs channel, wait for printer, then print footer
			close(jobs)
			WaitForPrinter()
			fmt.Printf("\n  ],\n")
			fmt.Printf("  \"last_block\": %d,\n", lastblock-1)
			fmt.Printf("  \"last_block_size\": %d,\n", lastread)
			fmt.Printf("  \"last_block_diff\": %d\n", Options.BlockSize-lastread)
			fmt.Printf("}\n")

		// --------------------------------------------------------------------
		} else {

			// Open connection to rpc
			//"/home/local/ANT/dieamare/dev/go/src/github.com/DiegoPomares/bigsync"
			//client, err := ClientRPC("bash", "-c", "go", "run", "*.go", "--server-mode-filename", "a", "dummy")
			client, err := ClientRPC("bash", "-c", "go run *.go --server-mode-filename a dummy")
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


func read_file(fh io.Reader, bs int, blocks chan<- Block) (int, int, error) {
	// Iterate through blocks in file
	var i, lastread int
	for i = 0;; i++ {
		buf := make([]byte, bs)

		init_time := time.Now()

		read_size, err := fh.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return i, 0, err
		}
		lastread = read_size

		blocks <- Block{i, buf[:read_size], []byte{}}

		if Verbose {
			u.Stderrln("Block", i, "read in", time.Now().Sub(init_time))
		}
	}

	return i, lastread, nil
}
