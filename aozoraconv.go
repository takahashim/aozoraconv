package aozoraconv

import (
	"strings"
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
	aozoraUtf8CharReplacerR = strings.NewReplacer(Reverse(aozoraCharMap)...)
)

// Reverse reverses aozoraUtf8CharReplacer
func Reverse(s []string) []string {
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
