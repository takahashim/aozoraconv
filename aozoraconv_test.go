package aozoraconv

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// toSjis convert string into Shift_JIS ([]byte)
func toSjis(s string) []byte {
	outBuf := new(bytes.Buffer)
	reader := strings.NewReader(s)
	writer := transform.NewWriter(outBuf, japanese.ShiftJIS.NewEncoder())

	//Copy
	_, err := io.Copy(writer, reader)
	if err != nil {
		return nil
	}
	ret := outBuf.Bytes()
	return ret
}

func TestAozoraConv(t *testing.T) {

	var convertedStrings = []struct {
		in  string
		out string
	}{
		{"", ""},
		{"あ", "あ"},
		{"〜", "～"},
		{"‖", "∥"},
		{"∥", "∥"},
		{"¢", "￠"},
	}
	for _, tt := range convertedStrings {
		if got, want := Conv(tt.in), tt.out; got != want {
			t.Errorf("Conv: got %v want %v", got, want)
		}
	}
}

func TestAozoraConvR(t *testing.T) {

	var convertedStrings = []struct {
		in  string
		out string
	}{
		{"", ""},
		{"あ", "あ"},
		{"～", "〜"},
		{"\u301C", "\u301C"},
		{"\uFF5E", "\u301C"},
		{"∥", "‖"},
		{"‖", "‖"},
		{"￠", "¢"},
	}
	for _, tt := range convertedStrings {
		if got, want := ConvRev(tt.in), tt.out; got != want {
			t.Errorf("ConvRev: got %v want %v", got, want)
		}
	}
}

func TestSjisConv(t *testing.T) {
	encoder := japanese.ShiftJIS.NewEncoder()

	input := "\uFF5E∥￠\u2015"
	reader := transform.NewReader(strings.NewReader(input), encoder)
	ret, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Errorf("sjis: error %v", err)
	}
	s := ""
	for _, c := range ret {
		s += fmt.Sprintf("%x", c)
	}
	want := "8160" + "8161" + "8191" + "815c"
	if s != want {
		t.Errorf("sjis: got %v want %v", s, want)
	}
}

func TestSjisConv2(t *testing.T) {
	input := "\uFF5E∥￠\u2015"
	out := new(bytes.Buffer)

	writer := transform.NewWriter(out, japanese.ShiftJIS.NewEncoder())
	n, err := writer.Write([]byte(input))
	if err != nil {
		t.Errorf("sjis: in column %v, val: %v, 'U+%x'", n, err, []rune(input[n:])[0])
	}
	got := ""
	for _, c := range out.Bytes() {
		got += fmt.Sprintf("%x", c)
	}
	want := "8160" + "8161" + "8191" + "815c"
	if got != want {
		t.Errorf("sjis: got %v want %v", got, want)
	}
}

func TestSjisConv3(t *testing.T) {
	input := "\uFF5E∥￠\u2015 \u2014123"
	out := new(bytes.Buffer)

	writer := transform.NewWriter(out, japanese.ShiftJIS.NewEncoder())
	n, err := writer.Write([]byte(input))
	if n != 13 || err == nil {
		t.Errorf("should be fail: %v, val: %v, 'U+%x'", n, err, []rune(input[n:])[0])
	}
}

func TestEncode(t *testing.T) {
	var convertedPairs = []struct {
		in  string
		out []byte
	}{
		{"あいうえお", toSjis("あいうえお")},
		{"\u301C", toSjis("\uFF5E")},
		{"\uFF5E", toSjis("\uFF5E")},
		{"¢", toSjis("￠")},
	}
	for _, tt := range convertedPairs {
		input := strings.NewReader(tt.in)
		output := new(bytes.Buffer)

		Encode(input, output)
		if got, want := output.Bytes(), tt.out; bytes.Compare(got, want) != 0 {
			t.Errorf("Encode got: %v, want: %v", got, want)
		}
	}
}

func TestDecode(t *testing.T) {
	var convertedPairs = []struct {
		out string
		in  []byte
	}{
		{"あいうえお", toSjis("あいうえお")},
		{"\u301C", toSjis("\uFF5E")},
		{"¢", toSjis("￠")},
	}
	for _, tt := range convertedPairs {
		input := bytes.NewBuffer(tt.in)
		output := new(bytes.Buffer)

		Decode(input, output)
		if got, want := output.String(), tt.out; got != want {
			t.Errorf("Encode got: %v, want: %v", got, want)
		}
	}
}
