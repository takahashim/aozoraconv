package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"io"

	"github.com/takahashim/aozoraconv"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func main() {
	const encSjis = 1
	const encUtf8 = 2

	var useSjis, useUtf8 bool
	var enc int
	var path, outpath string
	var input io.Reader
	var output io.Writer
	var err error

	flag.BoolVar(&useSjis, "sjis", false, "convert from UTF-8 into Shift_JIS")
	flag.BoolVar(&useUtf8, "utf8", false, "convert from Shift_JIS into UTF-8")
	flag.StringVar(&path, "f", "", "input filename")
	flag.StringVar(&outpath, "o", "", "output filename")

	flag.Parse()

	if path == "" {
		input = os.Stdin
	} else {
		input, err = os.Open(path)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}

	if outpath == "" {
		output = os.Stdout
	} else {
		output, err = os.OpenFile(outpath, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Fatal(err)
		}
	}

	if useUtf8 {
		enc = 2
	} else {
		enc = 1
	}

	if enc == encUtf8 {
		decoder := japanese.ShiftJIS.NewDecoder()
		reader := transform.NewReader(input, decoder)
		ret, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal(err)
		}
		str := aozoraconv.ConvRev(string(ret))
		_, err = fmt.Fprint(output, str)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		ret, err := ioutil.ReadAll(input)
		if err != nil {
			log.Fatal(err)
		}
		str := aozoraconv.Conv(string(ret))
		encoder := japanese.ShiftJIS.NewEncoder()
		writer := transform.NewWriter(output, encoder)
		_, err = fmt.Fprint(writer, str)
		if err != nil {
			log.Fatal(err)
		}
	}
}
