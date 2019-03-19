package match

import (
	"bufio"
	"bytes"
	"strconv"
	"testing"
)

var data = []byte("\nvalue\r\n" +
	"*3\r\n$3\r\nSET\r\n$1\r\na\r\n$1\r\n1\r\n" +
	"*4\r\n$4\r\nHSET\r\n$2\r\nxy\r\n$1\r\nz\r\n$1\r\n2\r\n" +
	"*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n" +
	"*4\r\n$4\r\nHSET\r")

func CreateSplitMatcher() *SplitMatcher {
	return NewSplitMatcher([]byte("*"), []byte("\r\n"))
}

func CreateFieldMatcher(chunk []byte) *FieldMatcher {
	// 举例：*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
	offset := bytes.IndexByte(chunk, byte('$'))
	index := bytes.Index(chunk[offset:], []byte("\r\n"))
	length, _ := strconv.Atoi(string(chunk[offset+1 : offset+index]))
	m := NewFieldMatcher()
	m.AddFixeds([]int{1, offset - 3, 2, 1, index - 1, 2, length, 2})
	return m
}

func SplitTestData(split bufio.SplitFunc) ([][]byte, error) {
	r := bytes.NewReader(data)
	return SplitStream(r, split)
}

func MatchChunk(chunk []byte, fm *FieldMatcher) (cmd string) {
	fs, _ := MatchHead(chunk, fm.Sequence)
	if len(fs) >= 7 {
		cmd = string(fs[6])
	}
	return
}

// 测试切割出完整的包
func TestSplit(t *testing.T) {
	tm := NewSplitMatcher([]byte("*"), []byte("*"))
	output, err := SplitTestData(tm.Spliter)
	if err != nil {
		t.Error(err)
	}
	for _, chunk := range output {
		t.Log(strconv.Quote(string(chunk)))
	}
}

// 测试切割出完整的包
func TestBetween(t *testing.T) {
	tm := CreateSplitMatcher()
	output, err := SplitTestData(tm.Spliter)
	if err != nil {
		t.Error(err)
	}
	for _, chunk := range output {
		t.Log(strconv.Quote(string(chunk)))
	}
}

// 测试切割出完整的包
func TestMatch(t *testing.T) {
	var fm *FieldMatcher
	tm := CreateSplitMatcher()
	output, _ := SplitTestData(tm.Spliter)
	for _, chunk := range output {
		if fm == nil {
			fm = CreateFieldMatcher(chunk)
		}
		cmd := MatchChunk(chunk, fm)
		t.Log(cmd)
	}
}

func BenchmarkSplit1(b *testing.B) {
	tm := CreateSplitMatcher()
	r := bytes.NewBuffer(data)
	for i := 0; i < b.N; i++ {
		SplitStream(r, tm.Spliter)
		r.Reset()
	}
}

func BenchmarkSplit2(b *testing.B) {
	tm := CreateSplitMatcher()
	r := bytes.NewReader(data)
	for i := 0; i < b.N; i++ {
		SplitStream(r, tm.Spliter)
	}
}

func BenchmarkSplit3(b *testing.B) {
	tm := CreateSplitMatcher()
	r := bytes.NewReader(data)
	for i := 0; i < b.N; i++ {
		SplitStream(r, tm.Spliter)
	}
}

func BenchmarkMatch(b *testing.B) {
	var fm *FieldMatcher
	tm := CreateSplitMatcher()
	output, _ := SplitTestData(tm.Spliter)
	size := len(output)
	if size > 0 {
		fm = CreateFieldMatcher(output[0])
	}
	for i := 0; i < b.N; i += size {
		MatchChunk(output[i%size], fm)
	}
}
