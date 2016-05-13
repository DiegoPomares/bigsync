package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const (
	NAME    = "bigsync"
	VERSION = "0.1 alpha"
	URL     = "https://github.com/DiegoPomares/bigsync"
)

var Verbose bool

var Options struct {
	SourceFile string
	RemoteHost string
	DestFile   string

	ServerMode bool

	BlockSize     int
	Workers       int
	HashType      string
	ForceCreation bool

	ExtraSsh string
	CustomSh string

	Diff  bool
	Equal bool
}

var DefaultBlockSize = (1024 * 1024)
var DefaultWorkers = 16 * runtime.NumCPU()
var DefaultHashType = "SHA256"

var usage = fmt.Sprintf(`%s %s (%s)
Usage: %s [options] source_file [ [[user@]remote_host] dest_file ]

Common:
    -b, --block-size <bytes>        Read and write up to <bytes> at a time.
                                    Supports multiplicative suffixes KMG
                                    (default: 1M = 1024*1024)
    -w, --workers <number>          Number of hashing workers (default: `+strconv.Itoa(DefaultWorkers)+`)
    -t, --hash-type <algoritm>      Available: MD5, SHA1, SHA256, SHA512 (default: `+DefaultHashType+`)
    -f, --force-creation            Create the file on the remote_host if it
                                    doesn't exist already

Transport:
    -s, --ssh <options>             Extra parameters for SSH
    -c, --custom-shell              Custom shell to communicate.
                                    Overrides SSH, and remote_host
    	
No-sync operations:
    -d, --diff                      Print different blocks list and exit
    -q, --equal                     Print equal blocks list and exit

Miscellaneous:
    -v, --verbose                   Print debugging messages
    -h, --help                      Show this message

See the %s manpage for full options, descriptions and usage examples
`, NAME, VERSION, URL, NAME, NAME+"(1)")

func print_usage() {
	fmt.Printf("%s", usage)
	os.Exit(0)
}

func arg_get(args []string, n int) string {
	if len(args) > n {
		return args[n]
	}

	return ""
}

func Process_opts() error {
	var block_size_human string

	// Flags
	flag.StringVar(&block_size_human, "b", "", "Block size")
	flag.StringVar(&block_size_human, "block-size", "", "Block size")
	flag.IntVar(&Options.Workers, "w", 0, "Workers")
	flag.IntVar(&Options.Workers, "workers", 0, "Workers")
	flag.StringVar(&Options.HashType, "t", "", "Hash type")
	flag.StringVar(&Options.HashType, "hash-type", "", "Hash type")
	flag.BoolVar(&Options.ForceCreation, "f", false, "Force creation")
	flag.BoolVar(&Options.ForceCreation, "force-creation", false, "Force creation")

	flag.StringVar(&Options.ExtraSsh, "s", "", "SSH extras")
	flag.StringVar(&Options.ExtraSsh, "ssh", "", "SSH extras")
	flag.StringVar(&Options.CustomSh, "c", "", "CustomSh")
	flag.StringVar(&Options.CustomSh, "custom-shell", "", "CustomSh")

	flag.BoolVar(&Options.Diff, "d", false, "Diff")
	flag.BoolVar(&Options.Diff, "diff", false, "Diff")
	flag.BoolVar(&Options.Equal, "q", false, "Equal")
	flag.BoolVar(&Options.Equal, "equal", false, "Equal")

	flag.BoolVar(&Verbose, "v", false, "Verbose")
	flag.BoolVar(&Verbose, "verbose", false, "Verbose")

	flag.BoolVar(&Options.ServerMode, "server-mode-enable", false, "Server mode")

	flag.Usage = print_usage
	flag.Parse()

	// Positional arguments
	Options.SourceFile = arg_get(flag.Args(), 0)

	Options.DestFile = arg_get(flag.Args(), 2)
	if Options.DestFile == "" {
		Options.DestFile = arg_get(flag.Args(), 1)	
	} else {
		Options.RemoteHost = arg_get(flag.Args(), 1)
	}


	// Check options
	block_size := DefaultBlockSize
	if block_size_human != "" {
		var err error
		block_size, err = parse_num(block_size_human)
		if err != nil {
			return err
		}

	}

	if Options.SourceFile == "" {
		return errors.New("Missing argument: source_file")
	}

	// Merge default opts
	Options.BlockSize = block_size
	if Options.Workers == 0 {
		Options.Workers = DefaultWorkers
	}
	if Options.HashType == "" {
		Options.HashType = DefaultHashType
	} else {
		Options.HashType = strings.ToUpper(Options.HashType)
	}

	return nil
}

func parse_num(s string) (int, error) {
	multi := 1
	s = strings.ToUpper(s)

	switch s[len(s)-1:] {
	case "K":
		multi = 1024
	case "M":
		multi = 1024 * 1024
	case "G":
		multi = 1024 * 1024 * 1024
	case "T":
		multi = 1024 * 1024 * 1024 * 1024
	default:
		s = s + " "
	}

	out, err := strconv.Atoi(s[:len(s)-1])

	return out * multi, err
}
