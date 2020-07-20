package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	var decode, dump, pretty bool
	flag.BoolVar(&decode, "d", false, "decodes input")
	flag.BoolVar(&dump, "c", false, "encodes the input as hexadecimal followed by characters")
	flag.BoolVar(&pretty, "p", false, "encoded using a prettier format aa:bb")
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
		if pretty {
			prettify(b)
		} else {
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

func prettify(data []byte) {
	last := len(data) - 1
	for i, b := range data {
		if i != 0 && (i%16) == 0 {
			fmt.Print("\n")
		}
		fmt.Printf("%02x", b)
		if i != last {
			fmt.Print(":")
		}
	}
	fmt.Println()
}
