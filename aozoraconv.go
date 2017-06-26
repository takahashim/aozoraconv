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

const (
	// EncSjis is magic number of Shift_JIS
	EncSjis = 1

	// EncUtf8 is magic number of UTF-8
	EncUtf8 = 2
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

// Jis2Uni returns a string from jis codepoint
func Jis2Uni(men, ku, ten int) (str string, err error) {
	if men < 1 || men > 2 || ku < 1 || ku > 94 || ten < 1 || ten > 94 {
		return "", fmt.Errorf("error: args should be in 1..2, 1..94, 1..94")
	}
	chr := jis0213Decode[men-1][ku-1][ten-1]
	if chr == "" {
		return "", fmt.Errorf("invalid access men: %v ku:%v ten:%v", men, ku, ten)
	}
	return chr, nil
}

// Uni2Jis returns a pointer of JisEntry
func Uni2Jis(str string) (jis JisEntry, err error) {
	r := []rune(str)
	r1 := r[0]
	if len(r) == 1 {
		switch {
		case 0x20 <= r1 && r1 < 0x7f:
			return JisEntry{0, 0, 0}, fmt.Errorf("ASCII character")
		case encode0Low <= r1 && r1 < encode0High:
			if r1 = rune(encode0[r1-encode0Low]); uint8(r1>>planeShift)&0x03 > 0 {
				goto write2
			}
		case encode1Low <= r1 && r1 < encode1High:
			if r1 = rune(encode1[r1-encode1Low]); uint8(r1>>planeShift)&0x03 > 0 {
				goto write2
			}
		case encode2Low <= r1 && r1 < encode2High:
			if r1 = rune(encode2[r1-encode2Low]); uint8(r1>>planeShift)&0x03 > 0 {
				goto write2
			}
		case encode3Low <= r1 && r1 < encode3High:
			if r1 = rune(encode3[r1-encode3Low]); uint8(r1>>planeShift)&0x03 > 0 {
				goto write2
			}
		case encode4Low <= r1 && r1 < encode4High:
			if r1 = rune(encode4[r1-encode4Low]); uint8(r1>>planeShift)&0x03 > 0 {
				goto write2
			}
		}
		return JisEntry{0, 0, 0}, fmt.Errorf("invalid character")
	write2:
		men := int8(r1 >> planeShift)
		ku := int8(uint8(r1>>codeShift) & codeMask)
		ten := int8(uint8(r1) & codeMask)
		return JisEntry{men: men, ku: ku, ten: ten}, nil
	} else if len(r) == 2 {
		r2 := r[1]
		entry, b := multichars[r1][r2]
		if !b {
			return JisEntry{0, 0, 0}, err
		}
		return entry, nil
	}
	return JisEntry{0, 0, 0}, fmt.Errorf("length of string should be 1 or 2")
}
