// Copyright 2017 Masayoshi Takahashi
//
// original copyright:
// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

// This program generates tables.go:
//	go run maketables.go | gofmt > /tmp/tables.go && cp /tmp/tables.go .

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"sort"
	"strings"

	"github.com/takahashim/aozoraconv"
)

//JisEntry is jis character with men, ku, ten
type JisEntry struct {
	men, ku, ten int
}

type UnicodeTbl [65536 * 4]JisEntry
type JisTbl [2][94][94]string

var reverse UnicodeTbl
var mapping JisTbl
var multichars map[int32]map[int32]JisEntry

func getTable(url string) {
	for i := range reverse {
		reverse[i].men = -1
		reverse[i].ku = -1
		reverse[i].ten = -1
	}
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Get: %v", err)
	}
	defer res.Body.Close()

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if s == "" || s[0] == '#' {
			continue
		}
		m, k, t, uni, uni2 := 0, 0, 0, int32(0), int32(0)

		if err = aozoraconv.ParseLine(s, &m, &k, &t, &uni, &uni2); err != nil {
			log.Fatalf("could not parse %q; %v", s, err)
		}
		m -= 2
		k -= 32
		t -= 32
		if m < 1 || 2 < m {
			log.Fatalf("JIS code men %d is out of range", m)
		}
		if k < 1 || 94 < k {
			log.Fatalf("JIS code ku %d is out of range", k)
		}
		if t < 1 || 94 < t {
			log.Fatalf("JIS code ten %d is out of range", t)
		}
		e := JisEntry{men: m, ku: k, ten: t}
		if uni2 > 0 {
			mapping[m-1][k-1][t-1] = string([]rune{uni, uni2})
			if _, exist := multichars[uni]; !exist {
				multichars[uni] = make(map[int32]JisEntry)
			}
			multichars[uni][uni2] = e
		} else if uni > 0 {
			mapping[m-1][k-1][t-1] = string([]rune{uni})
			if reverse[uni].men == -1 {
				reverse[uni] = e
			} else {
				log.Fatalf("%U is duplicated, %v %v and %v\n%q", uni, uni2, e, reverse[uni], s)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("scanner error: %v", err)
	}
}

func main() {
	fmt.Printf("// generated by go run maketables.go; DO NOT EDIT\n\n")
	fmt.Printf("// Package aozoraconv provides Aozora Bunko format encodings (JIS X 0208/Shift_JIS).\n")
	fmt.Printf(`package aozoraconv // import "github.com/takahashim/aozoraconv"` + "\n\n")

	reverse = UnicodeTbl{}
	mapping = JisTbl{}
	multichars = make(map[int32]map[int32]JisEntry)

	url := "http://x0213.org/codetable/jisx0213-2004-std.txt"
	getTable(url)

	fmt.Printf("//JisEntry is jis character with men, ku, ten\n")
	fmt.Printf("type JisEntry struct {\n	men, ku, ten int8\n}\n")

	fmt.Printf("// jis0213Decode is the decoding table from JIS 0213 code to Unicode.\n// It is defined at %s\n",
		url)
	fmt.Printf("var jis0213Decode = [2][94][94]string{\n")
	var counter int = 0
	for _, m1 := range mapping {
		fmt.Printf("\t{\n")
		for _, m2 := range m1 {
			fmt.Printf("\t\t{")
			counter = 0
			for _, m3 := range m2 {
				if m3 != "" {
					fmt.Printf("\t%q,", m3)
					counter++
					if counter >= 8 {
						counter = 0
						fmt.Printf("\n\t\t")
					}
				}
			}
			fmt.Printf("\t},\n")
		}
		fmt.Printf("\t},\n")
	}
	fmt.Printf("}\n\n")

	// Any run of at least separation continuous zero entries in the reverse map will
	// be a separate encode table.
	const separation = 1024

	intervals := []interval(nil)
	low, high := -1, -1
	for i, v := range reverse {
		if v.men < 1 {
			continue
		}
		if low < 0 {
			low = i
		} else if i-high >= separation {
			if high >= 0 {
				intervals = append(intervals, interval{low, high})
			}
			low = i
		}
		high = i + 1
	}
	if high >= 0 {
		intervals = append(intervals, interval{low, high})
	}
	sort.Sort(byDecreasingLength(intervals))

	fmt.Printf("const (\n")
	fmt.Printf("\tcodeMask   = 0x7f\n")
	fmt.Printf("\tcodeShift  = 7\n")
	fmt.Printf("\tplaneShift = 14\n")
	fmt.Printf(")\n\n")

	fmt.Printf("const numEncodeTables = %d\n\n", len(intervals))
	fmt.Printf("// encodeX are the encoding tables from Unicode to JIS code,\n")
	fmt.Printf("// sorted by decreasing length.\n")
	for i, v := range intervals {
		fmt.Printf("// encode%d: %5d entries for runes in [%6d, %6d).\n", i, v.len(), v.low, v.high)
	}
	fmt.Printf("//\n")
	fmt.Printf("// The high two bits of the value record whether the JIS code comes from the\n")
	fmt.Printf("// JIS X 0213:2004 plane 1 (high bits == 1) or JIS X 0213:2000 plane 2 (high bits == 2).\n")
	fmt.Printf("// The low 14 bits are two 7-bit unsigned integers j1 and j2 that form the\n")
	fmt.Printf("// JIS code (94*j1 + j2) within that table.\n")
	fmt.Printf("\n")

	for i, v := range intervals {
		fmt.Printf("const encode%dLow, encode%dHigh = 0x%x, 0x%x\n\n", i, i, v.low, v.high)
		fmt.Printf("var encode%d = [...]uint16{\n", i)
		for j := v.low; j < v.high; j++ {
			x := reverse[j]
			if x.men == -1 {
				continue
			}
			fmt.Printf("\t0x%x - 0x%x: %d<<14 | %2d<<7 | %2d, // %q\n",
				j, v.low, x.men, x.ku, x.ten, string(rune(j)))
		}
		fmt.Printf("}\n\n")
	}

	keys1 := reflect.ValueOf(multichars).MapKeys()
	sort.Slice(keys1, func(i, j int) bool {
		return keys1[i].Int() < keys1[j].Int()
	})

	fmt.Printf("var multichars = map[int32]map[int32]JisEntry{\n")
	for _, k1 := range keys1 {
		u1 := int32(k1.Int())
		fmt.Printf("\t0x%X: {\n", u1)
		keys2 := reflect.ValueOf(multichars[u1]).MapKeys()
		sort.Slice(keys2, func(i, j int) bool {
			return keys2[i].Int() < keys2[j].Int()
		})
		for _, k2 := range keys2 {
			u2 := int32(k2.Int())
			v := multichars[u1][u2]
			fmt.Printf("\t\t0x%X: JisEntry{men: %d, ku: %d, ten: %d},\n", u2, v.men, v.ku, v.ten)
		}
		fmt.Printf("\t},\n")
	}
	fmt.Printf("}\n")
}

// interval is a half-open interval [low, high).
type interval struct {
	low, high int
}

func (i interval) len() int { return i.high - i.low }

// byDecreasingLength sorts intervals by decreasing length.
type byDecreasingLength []interval

func (b byDecreasingLength) Len() int           { return len(b) }
func (b byDecreasingLength) Less(i, j int) bool { return b[i].len() > b[j].len() }
func (b byDecreasingLength) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
