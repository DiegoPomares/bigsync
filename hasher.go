package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/md5"
	"fmt"
	"sync"
	//"os"
	//"time"
	//u "github.com/DiegoPomares/bigsync/utils"
)

var WgHashWorkers sync.WaitGroup
var WgPrinter sync.WaitGroup

type Block struct {
	index int
	data  []byte
	size  int
}

type Hash struct {
	index int
	data  []byte
}

func wrap_sha1(data []byte) []byte {
	b_hash := sha1.Sum(data)
	return b_hash[:]
}

func wrap_sha256(data []byte) []byte {
	b_hash := sha256.Sum256(data)
	return b_hash[:]
}

func wrap_sha512(data []byte) []byte {
	b_hash := sha512.Sum512(data)
	return b_hash[:]
}

func wrap_md5(data []byte) []byte {
	b_hash := md5.Sum(data)
	return b_hash[:]
}

func get_hash_func(hash_alg string) (func([]byte) []byte) {
	switch hash_alg {
	case "SHA1": return wrap_sha1
	case "SHA256": return wrap_sha256
	case "SHA512": return wrap_sha512
	case "MD5": return wrap_md5
	}

	return nil
}

func hash_worker(jobs <-chan Block, results chan<- Hash, hash_func (func([]byte) []byte)) {

	for job := range jobs {

		hash := hash_func(job.data[:job.size])
		s_hash := hash[:]

		results <- Hash{job.index, s_hash}
	}
}

func print_worker(results <-chan Hash) {
	i := 0
	store := make(map[int]Hash)

	for result := range results {

		if result.index == i {
			print_result(result)
			i++

			for ; len(store) > 0; i++ {
				seg, ok := store[i]
				if !ok {
					break
				}

				print_result(seg)
				delete(store, i)
			}

		} else {
			store[result.index] = result
		}
	}
}

func print_result(result Hash) {

	if result.index != 0 {
		fmt.Printf(",\n")
	}
	fmt.Printf("    { \"block\": %d, \"hash\": \"%x\" }", result.index, result.data)

}

func WaitForPrinter() {
	WgPrinter.Wait()
}


func StartHashWorkers(n int, jobs <-chan Block, results chan<- Hash, hash_alg string) {
	WgHashWorkers.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer WgHashWorkers.Done()
			hash_worker(jobs, results, get_hash_func(hash_alg))
		}()
	}

	// Close results channel after workers are done
	go func() {
		WgHashWorkers.Wait()
		close(results)
	}()
}

func StartPrinter(results <-chan Hash) {
	WgPrinter.Add(1)

	go func () {
		defer WgPrinter.Done()
		print_worker(results)
	}()
}