package aozoraconv

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

var (
	aozoraCharMap = []string{
		"\u2014", "\u2015", // "—"
		"\u301C", "\uFF5E", // "〜"
		"\u2016", "\u2225", // "‖"
		"\u2212", "\uFF0D", // "−"
		"\u00A2", "\uFFE0", // "¢"
		"\u00A3", "\uFFE1", // "£"
		"\u00A5", "\uFFE5", // "¥"
		"\u00AC", "\uFFE2", // "¬"
	}
	aozoraUtf8CharReplacer  = strings.NewReplacer(aozoraCharMap...)
	aozoraUtf8CharReplacerR = strings.NewReplacer(reverse(aozoraCharMap)...)
)

// reverse reverses aozoraUtf8CharReplacer
func reverse(s []string) []string {
	r := make([]string, len(s))
	for i := len(r) - 1; i >= 0; i-- {
		opp := len(r) - i - 1
		r[i] = s[opp]
	}
	return r
}

// Conv replaces some characters in Unicode
func Conv(str string) string {
	return aozoraUtf8CharReplacer.Replace(str)
}

// ConvRev replaces some characters in Unicode
func ConvRev(str string) string {
	return aozoraUtf8CharReplacerR.Replace(str)
}

// Decode convert from UTF-8 into Aozora Bunko format (Shift_JIS)
func Decode(input io.Reader, output io.Writer) (err error) {
	decoder := japanese.ShiftJIS.NewDecoder()
	reader := transform.NewReader(input, decoder)
	ret, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	str := ConvRev(string(ret))
	_, err = fmt.Fprint(output, str)
	return err
}

// Encode convert from Aozora Bunko format (Shift_JIS) into UTF-8
func Encode(input io.Reader, output io.Writer) (err error) {
	ret, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}
	str := Conv(string(ret))
	encoder := japanese.ShiftJIS.NewEncoder()
	writer := transform.NewWriter(output, encoder)
	_, err = fmt.Fprint(writer, str)
	return err
}
