package aozoraconv

import (
	"testing"
)

func TestJis2013Decode(t *testing.T) {
	var convertedStrings = []struct {
		men, ku, ten int
		c            string
	}{
		{1, 1, 3, "。"},
		{1, 2, 1, "◆"},
		{1, 4, 2, "あ"},
		{1, 5, 2, "ア"},
		{1, 21, 21, "亀"},
		{2, 1, 1, "𠂉"},
	}
	for _, tt := range convertedStrings {
		if got, want := jis0213Decode[tt.men-1][tt.ku-1][tt.ten-1], tt.c; got != want {
			t.Errorf("jis0213Decode: got %v want %v", got, want)
		}
	}
}

func TestEncode1(t *testing.T) {
	var codes = []struct {
		in           uint16
		men, ku, ten int
	}{
		{0x3000, 1, 1, 1},
		{0x3006, 1, 1, 26},
		{0x3042, 1, 4, 2},
	}
	for _, tt := range codes {
		got, want := encode1[tt.in-encode1Low], uint16(tt.men<<14|tt.ku<<7|tt.ten)
		if got != want {
			t.Errorf("encode1: got %v want %v", got, want)
		}
	}
}
