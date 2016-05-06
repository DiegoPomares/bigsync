package main

import (
	"fmt"
	"os"
	"flag"
	"runtime"
	"strings"
	"strconv"
	"errors"
)

const (
	NAME = "bigsync"
	VERSION = "0.1 beta"
	URL = "https://github.com/DiegoPomares/bigsync"

	BLOCK_SIZE = (1024 * 1024)
	THREAD_MULTIPLIER = 16
	HASH_ALGORITHM = "SHA256"
)

var Options struct {
	SourceFile string
	User string
	RemoteHost string
	DestFile string

	ServerMode string

	BlockSize int
	Workers int
	HashAlgorithm string
	ForceCreation bool

	Diff bool
	Equal bool
	JustPrint bool

	Verbose bool
}

var usage = fmt.Sprintf(`%s %s (%s)
Usage: %s [options] source_file [user@]remote_host [dest_file]

Common:
    -b, --block-size <bytes>        Read and write up to <bytes> at a time.
                                    Supports multiplicative suffixes KMG
                                    (default: 1M = 1024*1024)
    -w, --workers <number>          Number of hashing workers (default: ` + strconv.Itoa(runtime.NumCPU() * THREAD_MULTIPLIER) + `)
    -a, --algoritm <algoritm>       Hashing algoritm (default: ` + HASH_ALGORITHM + `), use -l
                                    Available: MD5, SHA1, SHA256, SHA512
    --force-creation                Create the file on the remote_host if it
                                    doesn't exist already
    	
No-sync operations:
    -d, --diff                      Print different blocks list and exit
    -e, --equal                     Print equal blocks list and exit
    -j, --just-print                Print the hash list of the source_file
                                    and exit. remote_host is not required

Miscellaneous:
    -v, --verbose                   Print debugging messages
    -h, --help                      Show this message

See the %s manpage for full options, descriptions and usage examples
`, NAME, VERSION, URL, NAME, NAME + "(1)")

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

func parse_num(s string) (int, error) {
	multi := 1
	s = strings.ToUpper(s)

	switch s[len(s)-1:] {
	case "K":
		multi = 1024
	case "M":
		multi = 1024*1024
	case "G":
		multi = 1024*1024*1024
	case "T":
		multi = 1024*1024*1024*1024
	default:
		s = s + " "
	}

	out, err := strconv.Atoi(s[:len(s)-1])

	return out * multi, err
}

func Process_opts() error {
	var block_size_human string

	// Flags
	flag.StringVar(&block_size_human, "b", "", "Block size")
	flag.StringVar(&block_size_human, "block_size", "", "Block size")
	flag.IntVar(&Options.Workers, "w", 0, "Workers")
	flag.IntVar(&Options.Workers, "workers", 0, "Workers")
	flag.StringVar(&Options.HashAlgorithm, "a", "", "Hash algoritm")
	flag.StringVar(&Options.HashAlgorithm, "algoritm", "", "Hash algoritm")
	flag.BoolVar(&Options.ForceCreation, "force_creation", false, "Force creation")

	flag.BoolVar(&Options.Diff, "d", false, "Diff")
	flag.BoolVar(&Options.Diff, "diff", false, "Diff")
	flag.BoolVar(&Options.Equal, "e", false, "Equal")
	flag.BoolVar(&Options.Equal, "equal", false, "Equal")
	flag.BoolVar(&Options.JustPrint, "j", false, "JustPrint")
	flag.BoolVar(&Options.JustPrint, "just_print", false, "JustPrint")

	flag.BoolVar(&Options.Verbose, "v", false, "Verbose")
	flag.BoolVar(&Options.Verbose, "verbose", false, "Verbose")

	flag.StringVar(&Options.ServerMode, "server_mode_filename", "", "ServerMode")

	flag.Usage = print_usage
	flag.Parse()


	// Positional arguments
	Options.SourceFile = arg_get(flag.Args(), 0)
	Options.DestFile = arg_get(flag.Args(), 2)

	topt := arg_get(flag.Args(), 1)
	topts := strings.Split(topt, "@")
	if len(topts) == 2 {
		Options.User = topts[0]
		Options.RemoteHost = topts[1]
	} else {
		Options.RemoteHost = topts[0]
	}

	// Check options
	block_size := BLOCK_SIZE
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
		Options.Workers = runtime.NumCPU() * THREAD_MULTIPLIER
	}
	if Options.HashAlgorithm == "" {
		Options.HashAlgorithm = HASH_ALGORITHM
	} else {
		Options.HashAlgorithm = strings.ToUpper(Options.HashAlgorithm)
	}

return nil
}