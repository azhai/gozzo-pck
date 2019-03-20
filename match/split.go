package match

import (
	"bufio"
	"bytes"
	"io"
)

// 判断匹配开头或结尾的方法
type MatchFunc func(data []byte) (int, int)

func MatchTwice(match MatchFunc, data []byte, atEOF bool) (int, int, int) {
	var (
		advance = -1
		i, m    int
	)
	if atEOF && len(data) == 0 {
		return advance, 0, 0 // 终止
	}
	if i, m = match(data); i >= 0 {
		advance = i
		if j, _ := match(data[i+m:]); j >= 0 {
			advance += j + m
		}
	}
	return advance, i, m
}

// 根据相同的开头和结尾分割
func SplitBoth(matchBoth MatchFunc) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (int, []byte, error) {
		advance, i, m := MatchTwice(matchBoth, data, atEOF)
		if advance < 0 {
			return 0, nil, nil // 终止(atEOF=true)或请求更多数据
		}
		var token []byte
		if advance > i {
			advance += m
			token = data[i:advance]
		}
		return advance, token, nil
	}
}

// 根据不同的开头和结尾分割
func SplitBetween(matchStart, matchEnd MatchFunc) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (int, []byte, error) {
		advance, i, m := MatchTwice(matchStart, data, atEOF)
		if advance < 0 {
			return 0, nil, nil // 终止(atEOF=true)或请求更多数据
		}
		var token []byte
		if advance > i {
			token = data[i:advance]
		} else {
			advance = len(data)
			token = data[i:]
		}
		if j, n := matchEnd(token[m:]); j >= 0 {
			token = token[:m+j+n]
		} else {
			token = nil
		}
		return advance, token, nil
	}
}

// 按前后标记拆包
type SplitMatcher struct {
	StartToken []byte
	EndToken   []byte
	Spliter    bufio.SplitFunc
}

func NewSplitMatcher(start, end []byte) *SplitMatcher {
	return &SplitMatcher{StartToken: start, EndToken: end}
}

func (m *SplitMatcher) GetSpliter() bufio.SplitFunc {
	if m.Spliter == nil {
		MatchStart := func(data []byte) (int, int) {
			return bytes.Index(data, m.StartToken), len(m.StartToken)
		}
		if bytes.Compare(m.StartToken, m.EndToken) == 0 { // 相同
			m.Spliter = SplitBoth(MatchStart)
		} else {
			MatchEnd := func(data []byte) (int, int) {
				return bytes.LastIndex(data, m.EndToken), len(m.EndToken)
			}
			m.Spliter = SplitBetween(MatchStart, MatchEnd)
		}
	}
	return m.Spliter
}

//解析字节流
func (m *SplitMatcher) SplitStream(outch chan<- []byte, rd io.Reader) error {
	scanner := bufio.NewScanner(rd)
	scanner.Split(m.GetSpliter())
	for scanner.Scan() {
		if chunk := scanner.Bytes(); chunk != nil {
			outch <- chunk
		}
	}
	return scanner.Err()
}

//解析二进制数据
func (m *SplitMatcher) SplitBuffer(input []byte) ([][]byte, error) {
	var output [][]byte
	rd := bytes.NewReader(input)
	scanner := bufio.NewScanner(rd)
	scanner.Split(m.GetSpliter())
	for scanner.Scan() {
		if chunk := scanner.Bytes(); chunk != nil {
			output = append(output, chunk)
		}
	}
	return output, scanner.Err()
}
