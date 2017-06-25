package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/takahashim/aozoraconv"
)

func errorf(format string, a ...interface{}) (ret int, err error) {
	ret, err = fmt.Fprintf(os.Stderr, format+"\n", a...)
	return ret, err
}

func getOuput(path string) (output io.Writer, err error) {
	if path == "" {
		return os.Stdout, nil
	}
	output, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func getInput(path string) (input io.Reader, err error) {
	if path == "" {
		return nil, errors.New("input file is not defined")
	}
	input, err = os.Open(path)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func doMain() int {

	var (
		useSjis, useUtf8 bool
		enc              int
		path, outpath    string
		encoding         string
	)

	flag.StringVar(&encoding, "e", "sjis", "set output encoding (sjis or utf8)")
	flag.BoolVar(&useSjis, "s", false, "convert from UTF-8 into Shift_JIS")
	flag.BoolVar(&useUtf8, "u", false, "convert from Shift_JIS into UTF-8")
	flag.StringVar(&outpath, "o", "", "output filename")

	flag.Parse()

	path = flag.Arg(0)

	input, err := getInput(path)
	if err != nil {
		errorf("error: %s", err)
		return 1
	}

	output, err := getOuput(outpath)
	if err != nil {
		errorf("error: %s", err)
		return 1
	}

	if strings.ToLower(encoding) == "utf8" || strings.ToLower(encoding) == "utf-8" || useUtf8 {
		enc = aozoraconv.EncUtf8
	} else if strings.ToLower(encoding) == "sjis" || strings.ToLower(encoding) == "shift_jis" || useSjis {
		enc = aozoraconv.EncSjis
	} else {
		errorf("define encoding -s (Shift_JIS) or -u (UTF-8) or -e sting")
		return 1
	}

	if enc == aozoraconv.EncUtf8 {
		err = aozoraconv.Decode(input, output)
	} else { // enc == aozoraconv.EncSjis
		err = aozoraconv.Encode(input, output)
	}
	if err != nil {
		errorf("error: %v", err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(doMain())
}
