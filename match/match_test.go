package match

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var data = []byte("\nvalue\r\n" +
	"*3\r\n$3\r\nSET\r\n$1\r\na\r\n$1\r\n1\r\n" +
	"*4\r\n$4\r\nHSET\r\n$2\r\nxy\r\n$1\r\nz\r\n$1\r\n2\r\n" +
	"*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r")

func CreateFieldMatcher(chunk []byte) *FieldMatcher {
	// 举例：*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
	offset := bytes.IndexByte(chunk, byte('$'))
	index := bytes.Index(chunk[offset:], []byte("\r\n"))
	size, _ := strconv.Atoi(string(chunk[offset+1 : offset+index]))
	m := NewFieldMatcher()
	m.AddFixeds([]int{1, offset - 3, 2, 1, index - 1, 2}, nil)
	m.AddField("cmd", NewField(size, false))
	return m
}

func MatchChunk(chunk []byte, fm *FieldMatcher) (cmd string) {
	data := fm.Match(chunk, true)
	if len(data) >= 7 {
		cmd = string(data["cmd"])
	}
	return
}

// 测试切割出完整的包
func TestInclusio(t *testing.T) {
	outch := make(chan []byte)
	go func() {
		for chunk := range outch {
			assert.Equal(t, byte('*'), chunk[0])
			tail := chunk[len(chunk)-1]
			assert.Equal(t, byte('*'), tail)
			t.Log(strconv.Quote(string(chunk)))
		}
	}()
	sp := NewSplitMatcher([]byte("*"), []byte("*"))
	err := sp.SplitStream(outch, bytes.NewReader(data))
	assert.NoError(t, err)
}

// 测试切割出完整的包
func TestBetween(t *testing.T) {
	sp := NewSplitMatcher([]byte("*"), []byte("\r\n"))
	output, err := sp.SplitBuffer(data)
	assert.NoError(t, err)
	for _, chunk := range output {
		assert.Equal(t, byte('*'), chunk[0])
		tail := chunk[len(chunk)-2:]
		assert.Equal(t, []byte("\r\n"), tail)
		t.Log(strconv.Quote(string(chunk)))
	}
}

// 测试切割出完整的包
func TestMatch(t *testing.T) {
	sp := NewSplitMatcher([]byte("*"), []byte("\r\n"))
	output, err := sp.SplitBuffer(data)
	assert.NoError(t, err)
	assert.Len(t, output, 3)
	for i, chunk := range output {
		fm := CreateFieldMatcher(chunk)
		cmd := MatchChunk(chunk, fm)
		if i == 1 {
			assert.Equal(t, "HSET", cmd)
		} else {
			assert.Equal(t, "SET", cmd)
		}
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
	assert.NoError(b, err)
	for _, chunk := range output {
		fms = append(fms, CreateFieldMatcher(chunk))
	}
	for i := 0; i < b.N; i += len(fms) {
		for j, chunk := range output {
			MatchChunk(chunk, fms[j])
		}
	}
}
