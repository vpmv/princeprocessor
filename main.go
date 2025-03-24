package main

import (
	"bufio"
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
	"strings"
)

var (
	minElems   *int
	maxElems   *int
	wordlist   *string
	outputFile *string
)

type outFile struct {
	fp  *os.File
	buf *bufio.Writer
}

// Reads the wordlist into a slice of words
func processWordlist(inputFile string) ([]string, error) {
	var scanner *bufio.Scanner
	var words []string

	if inputFile == "" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		file, err := os.Open(inputFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	}

	// reading words from input
	fmt.Println("No wordlist detected.\nYou can now manually enter words; Please submit your first word:")
	var lineBreaks int
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			words = append(words, line)
		} else {
			lineBreaks++
			if lineBreaks >= 2 {
				break
			}
			fmt.Println("Press [ENTER] one more time to submit your wordlist.")
		}
	}

	return words, scanner.Err()
}

func init() {
	minElems = flag.IntP("min", "n", 2, "Minimum number of elements per chain")
	maxElems = flag.IntP("max", "m", 4, "Maximum number of elements per chain")
	wordlist = flag.StringP("wordlist", "i", "", "Path to input wordlist file. Use stdin when omitted")
	outputFile = flag.StringP("output", "o", "", "Output file. Use stdout when omitted")
	wordSeparator := flag.StringP("separator", "s", string(Separator), "Separator used between elements")
	flag.Parse()

	if len(os.Args) > 1 && os.Args[1] == `help` {
		flag.Usage()
		os.Exit(0)
	}

	// Validate element count range
	if *minElems < 1 || *maxElems < *minElems {
		fmt.Printf("Error: --min must be ≥ 1 and ≤ --max\n")
		os.Exit(1)
	}

	// Read words from the wordlist
	var err error
	Dictionary, err = processWordlist(*wordlist)
	if err != nil {
		fmt.Printf("Failed to process wordlist: %v\n", err)
		os.Exit(1)
	}
	Separator = []byte(*wordSeparator)[0]
}

func main() {
	var (
		fp  *os.File
		err error
	)

	if *outputFile != "" {
		fp, err = os.Create(*outputFile)
		if err != nil {
			fmt.Printf("Error opening output file: %v\n", err)
			os.Exit(1)
		}
	} else {
		fp = os.Stdout
	}

	defer fp.Close()

	output := &outFile{fp, bufio.NewWriterSize(fp, 16*1024*1024)}

	generateChains(*minElems, *maxElems, output)
	output.buf.Flush()
}
