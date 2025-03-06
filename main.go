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
	"strconv"
	"strings"
)

type prettyFormat struct {
	IsSet   bool
	Columns int
}

func (p *prettyFormat) String() string {
	if p.Columns == 0 {
		return "true"
	}
	return strconv.Itoa(p.Columns)
}

func (p *prettyFormat) Set(s string) error {
	s = strings.ToLower(s)
	if s == "true" || s == "false" {
		p.IsSet = (s == "true")
		return nil
	}

	n, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	p.IsSet = true
	p.Columns = n
	return nil
}

func (p *prettyFormat) IsBoolFlag() bool {
	return true
}

func main() {
	var decode, dump, number bool
	var base int
	var pretty, goFormat prettyFormat

	flag.BoolVar(&decode, "d", false, "decodes input")
	flag.BoolVar(&dump, "c", false, "encodes the input as hexadecimal followed by characters")
	flag.Var(&pretty, "p", "encoded using a prettier format aa:bb, pass -p=n to print using n columns")
	flag.Var(&goFormat, "go", "encoded using as Go's []byte, pass -go=n to print using n columns")
	flag.BoolVar(&number, "n", false, "interprets input as a number")
	flag.IntVar(&base, "b", 10, "base used to when -n is used")
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
		case pretty.IsSet:
			prettify(b, pretty.Columns)
		case goFormat.IsSet:
			goify(b, goFormat.Columns)
		default:
			fmt.Printf("%x\n", b)
		}
	}
}

func usage() {
	o := flag.CommandLine.Output()
	fmt.Fprintf(o, "Usage: %s [<filename>]\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(o, "  -c      %s\n", flag.Lookup("c").Usage)
	fmt.Fprintf(o, "  -d      %s\n", flag.Lookup("d").Usage)
	fmt.Fprintf(o, "  -p      %s\n", flag.Lookup("p").Usage)
	fmt.Fprintf(o, "  -go     %s\n", flag.Lookup("go").Usage)
	fmt.Fprintf(o, "  -n      %s\n", flag.Lookup("n").Usage)
	fmt.Fprintf(o, "  -b int  %s (default: %s)\n", flag.Lookup("b").Usage, flag.Lookup("b").DefValue)
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
