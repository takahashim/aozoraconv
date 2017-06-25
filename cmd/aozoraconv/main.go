package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/takahashim/aozoraconv"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func main() {
	os.Exit(doMain())
}

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
		return os.Stdin, nil
	}
	input, err = os.Open(path)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func doMain() int {
	const encSjis = 1
	const encUtf8 = 2

	var useSjis, useUtf8 bool
	var enc int
	var path, outpath string
	var encoding string
	var input io.Reader
	var output io.Writer
	var err error

	flag.StringVar(&encoding, "e", "sjis", "set output encoding (sjis or utf8)")
	flag.BoolVar(&useSjis, "s", false, "convert from UTF-8 into Shift_JIS")
	flag.BoolVar(&useUtf8, "u", false, "convert from Shift_JIS into UTF-8")
	flag.StringVar(&path, "f", "", "input filename")
	flag.StringVar(&outpath, "o", "", "output filename")

	flag.Parse()

	input, err = getInput(path)
	if err != nil {
		errorf("error: %v", err)
		return 1
	}

	output, err = getOuput(outpath)
	if err != nil {
		errorf("%s", err)
		return 1
	}

	if strings.ToLower(encoding) == "utf8" || strings.ToLower(encoding) == "utf-8" || useUtf8 {
		enc = 2
	} else if strings.ToLower(encoding) == "sjis" || strings.ToLower(encoding) == "shift_jis" || useSjis {
		enc = 1
	} else {
		errorf("define encoding -s (Shift_JIS) or -u (UTF-8) or -e sting")
		return 1
	}

	if enc == encUtf8 {
		decoder := japanese.ShiftJIS.NewDecoder()
		reader := transform.NewReader(input, decoder)
		ret, err := ioutil.ReadAll(reader)
		if err != nil {
			errorf("error: %v", err)
			return 1
		}
		str := aozoraconv.ConvRev(string(ret))
		_, err = fmt.Fprint(output, str)
		if err != nil {
			errorf("error: %v", err)
			return 1
		}
	} else {
		ret, err := ioutil.ReadAll(input)
		if err != nil {
			errorf("error: %v", err)
			return 1
		}
		str := aozoraconv.Conv(string(ret))
		encoder := japanese.ShiftJIS.NewEncoder()
		writer := transform.NewWriter(output, encoder)
		_, err = fmt.Fprint(writer, str)
		if err != nil {
			errorf("error: %v", err)
			return 1
		}
	}
	return 0
}
