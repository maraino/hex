package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
)

func main() {
	var decode, dump, pretty, number, goFormat bool
	var base, cols int
	flag.BoolVar(&decode, "d", false, "decodes input")
	flag.BoolVar(&dump, "c", false, "encodes the input as hexadecimal followed by characters")
	flag.BoolVar(&pretty, "p", false, "encoded using a prettier format aa:bb")
	flag.BoolVar(&goFormat, "go", false, "encoded using as Go's []byte")
	flag.BoolVar(&number, "n", false, "interprets input as a number")
	flag.IntVar(&base, "b", 10, "base used to when -n is used")
	flag.IntVar(&cols, "cols", 0, "number of columns for pretty and Go's format")
	flag.Usage = usage
	flag.Parse()

	var b []byte
	var err error
	switch len(flag.Args()) {
	case 0:
		b, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			fatal(err)
		}
	case 1:
		b, err = ioutil.ReadFile(flag.Arg(0))
		if err != nil {
			fatal(err)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}

	switch {
	case decode:
		var bb bytes.Buffer
		// Remove 0x prefix
		if len(b) >= 2 && bytes.EqualFold(b[:2], []byte("0x")) {
			b = b[2:]
		}
		for _, c := range b {
			if isHexChar(c) {
				bb.WriteByte(c)
			}
		}
		out := make([]byte, bb.Len()/2)
		if _, err = hex.Decode(out, bb.Bytes()); err != nil {
			fatal(err)
		}
		if dump {
			fmt.Print(hex.Dump(out))
		} else {
			os.Stdout.Write(out)
		}
	case dump:
		fmt.Print(hex.Dump(b))
	default:
		if number {
			n, ok := new(big.Int).SetString(string(bytes.TrimSpace(b)), base)
			if !ok {
				fatal(fmt.Errorf("error parsing %s in base %d", b, base))
			}
			b = n.Bytes()
		}
		switch {
		case pretty:
			prettify(b, cols)
		case goFormat:
			goify(b, cols)
		default:
			fmt.Printf("%x\n", b)
		}
	}
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [<filename>]\n", filepath.Base(os.Args[0]))
	flag.PrintDefaults()
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func isHexChar(c byte) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	default:
		return false
	}
}

func prettify(data []byte, cols int) {
	if cols == 0 {
		cols = 16
	}
	last := len(data) - 1
	for i, b := range data {
		if i != 0 && (i%cols) == 0 {
			fmt.Print("\n")
		}
		fmt.Printf("%02x", b)
		if i != last {
			fmt.Print(":")
		}
	}
	fmt.Println()
}

func goify(data []byte, cols int) {
	if cols == 0 {
		cols = 8
	}
	last := len(data) - 1
	fmt.Println("[]byte{")
	for i, b := range data {
		if (i % cols) == 0 {
			if i == 0 {
				fmt.Print("\t")
			} else {
				fmt.Print("\n\t")
			}
		}

		fmt.Printf("0x%02x", b)
		if i != last {
			fmt.Print(", ")
		} else {
			fmt.Print(",")
		}
	}
	fmt.Println("\n}")
}
