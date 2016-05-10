package hasher

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	//"fmt"
	"os"
	"sync"
	//"io"
	"errors"
	"math"
	//u "github.com/DiegoPomares/bigsync/utils"
)

const (
	READ_BUFFER = 4
)

type Block struct {
	Index int
	Data  []byte
	Size  int
	Hash  []byte
}

func New(file_path string, mode string, bs int, hash_type string, workers int) (*hasher, error) {
	var err error

	var file_flags int
	switch mode {
	case "r":
		file_flags = os.O_RDONLY
	case "rw":
		file_flags = os.O_RDWR //TODO check if O_CREATE is needed
	default:
		return nil, os.ErrInvalid
	}

	h := new(hasher)

	// Open file
	h.fh, err = os.OpenFile(file_path, file_flags, 0)
	if err != nil {
		return nil, err
	}

	// Get file info
	h.FileInfo, err = h.fh.Stat()
	if err != nil {
		return nil, err
	}

	switch hash_type {
	case "SHA1":
		h.hash_func = w_sha1
	case "SHA256":
		h.hash_func = w_sha256
	case "SHA512":
		h.hash_func = w_sha512
	case "MD5":
		h.hash_func = w_md5
	default:
		return nil, errors.New("invalid hash type")
	}

	// Initialize
	h.FilePath = file_path
	h.BlockSize = bs
	h.HashType = hash_type
	h.Workers = workers
	h.FileSize = h.FileInfo.Size()
	h.NumBlocks = int(math.Ceil(float64(h.FileSize) / float64(bs)))
	h.blocks = make(chan Block, READ_BUFFER)
	h.Hashes = make(chan Block, workers)

	return h, nil
}

type hasher struct {
	FileInfo  os.FileInfo
	FilePath  string
	BlockSize int
	Workers   int
	HashType  string
	FileSize  int64
	NumBlocks int

	fh        *os.File
	hash_func func([]byte) []byte
	blocks    chan Block
	Hashes    chan Block

	wg_workers sync.WaitGroup
}

func (self *hasher) Start() {

	self.start_workers(self.Workers)

	// Read the blocks in the file
	go func() {
		for i := 0; ; i++ {
			buf := make([]byte, self.BlockSize)

			read_size, err := self.fh.Read(buf)
			if err != nil {
				break
			}

			self.blocks <- Block{i, buf[:read_size], read_size, []byte{}}
		}

		close(self.blocks)
	}()

}

func (self *hasher) HashRange() chan Block {
	ch := make(chan Block, READ_BUFFER)
	
	go func() {
		current := 0
		store := make(map[int]Block)

		for result := range self.Hashes {

			if result.Index == current {
				ch <- result
				current++

				for ; len(store) > 0; current++ {
					block, ok := store[current]
					if !ok {
						break
					}

					ch <- block
					delete(store, current)
				}

			} else {
				store[result.Index] = result
			}
		}

		close(ch)
	}()

	return ch
}

func (self *hasher) start_workers(n int) {
	self.wg_workers.Add(n)

	// Spawn workers
	for i := 0; i < n; i++ {
		go func() {
			defer self.wg_workers.Done()
			self.hash_worker()
		}()
	}

	// Close Hashes channel after workers are done
	go func() {
		self.wg_workers.Wait()
		close(self.Hashes)
	}()
}

func (self *hasher) hash_worker() {
	for job := range self.blocks {

		hash := self.hash_func(job.Data)
		job.Hash = hash[:]

		self.Hashes <- job
	}
}


// Wrappers -------------------------------------------------------------------
func w_sha1(data []byte) []byte {
	b_hash := sha1.Sum(data)
	return b_hash[:]
}

func w_sha256(data []byte) []byte {
	b_hash := sha256.Sum256(data)
	return b_hash[:]
}

func w_sha512(data []byte) []byte {
	b_hash := sha512.Sum512(data)
	return b_hash[:]
}

func w_md5(data []byte) []byte {
	b_hash := md5.Sum(data)
	return b_hash[:]
}
