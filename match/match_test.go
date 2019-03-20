package match

import (
	"bytes"
	"strconv"
	"testing"
)

var data = []byte("\nvalue\r\n" +
	"*3\r\n$3\r\nSET\r\n$1\r\na\r\n$1\r\n1\r\n" +
	"*4\r\n$4\r\nHSET\r\n$2\r\nxy\r\n$1\r\nz\r\n$1\r\n2\r\n" +
	"*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r")

func CreateFieldMatcher(chunk []byte) *FieldMatcher {
	// 举例：*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
	offset := bytes.IndexByte(chunk, byte('$'))
	index := bytes.Index(chunk[offset:], []byte("\r\n"))
	length, _ := strconv.Atoi(string(chunk[offset+1 : offset+index]))
	m := NewFieldMatcher()
	m.AddFixeds([]int{1, offset - 3, 2, 1, index - 1, 2, length, 2})
	return m
}

func MatchChunk(chunk []byte, fm *FieldMatcher) (cmd string) {
	fs, _ := MatchHead(chunk, fm.Sequence)
	if len(fs) >= 7 {
		cmd = string(fs[6])
	}
	return
}

// 测试切割出完整的包
func TestInclusio(t *testing.T) {
	outch := make(chan []byte)
	go func() {
		for chunk := range outch {
			t.Log(strconv.Quote(string(chunk)))
		}
	}()
	sp := NewSplitMatcher([]byte("*"), []byte("*"))
	err := sp.SplitStream(outch, bytes.NewReader(data))
	if err != nil {
		t.Error(err)
	}
}

// 测试切割出完整的包
func TestBetween(t *testing.T) {
	sp := NewSplitMatcher([]byte("*"), []byte("\r\n"))
	output, err := sp.SplitBuffer(data)
	if err != nil {
		t.Error(err)
		return
	}
	for _, chunk := range output {
		t.Log(strconv.Quote(string(chunk)))
	}
}

// 测试切割出完整的包
func TestMatch(t *testing.T) {
	sp := NewSplitMatcher([]byte("*"), []byte("\r\n"))
	output, err := sp.SplitBuffer(data)
	if err != nil {
		t.Error(err)
		return
	}
	for _, chunk := range output {
		fm := CreateFieldMatcher(chunk)
		cmd := MatchChunk(chunk, fm)
		t.Log(cmd)
	}
}

func BenchmarkSplit(b *testing.B) {
	sp := NewSplitMatcher([]byte("*"), []byte("\r\n"))
	for i := 0; i < b.N; i++ {
		sp.SplitBuffer(data)
	}
}

func BenchmarkMatch(b *testing.B) {
	var fms []*FieldMatcher
	sp := NewSplitMatcher([]byte("*"), []byte("\r\n"))
	output, err := sp.SplitBuffer(data)
	if err != nil {
		b.Error(err)
		return
	}
	for _, chunk := range output {
		fms = append(fms, CreateFieldMatcher(chunk))
	}
	for i := 0; i < b.N; i += len(fms) {
		for j, chunk := range output {
			MatchChunk(chunk, fms[j])
		}
	}
}
