package match

import (
	"bufio"
	"bytes"
	"io"
)

// 判断匹配开头或结尾的方法
type MatchBytesFunc func(data []byte) (int, int)

// 根据开头几个字节分割
func SplitByHead(matchHead MatchBytesFunc, sameTail bool) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF {
			if advance = len(data); advance == 0 {
				return 0, nil, nil // 终止
			}
		}
		if i, n := matchHead(data); i >= 0 {
			if j, _ := matchHead(data[i+n:]); j >= 0 {
				// finds >= 2
				advance = i + n + j
				if sameTail {
					advance += n
				}
				token = data[i:advance]
				return advance, token, nil
			} else if sameTail == false {
				// finds = 1
				if atEOF {
					token = data[i:]
				} else {
					advance = i
				}
				return advance, token, nil
			}
		}
		// finds = 0
		return advance, nil, nil
	}
}

// 根据开头和结尾几个字节分割
func SplitBetween(matchHead, matchTail MatchBytesFunc) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		split := SplitByHead(matchHead, false)
		advance, token, err = split(data, atEOF)
		if len(token) > 0 {
			if i, n := matchTail(token); i >= 0 {
				token = token[:i+n]
			} else {
				token = nil
			}
		}
		return
	}
}

// 按前后标记拆包
type SplitMatcher struct {
	MatchStart MatchBytesFunc
	MatchEnd   MatchBytesFunc
	Spliter    bufio.SplitFunc
}

func NewSplitMatcher(start, end []byte) *SplitMatcher {
	m := new(SplitMatcher)
	m.MatchStart = func(data []byte) (int, int) {
		return bytes.Index(data, start), len(start)
	}
	if bytes.Compare(start, end) == 0 { // 相同
		m.Spliter = SplitByHead(m.MatchStart, true)
	} else {
		m.MatchEnd = func(data []byte) (int, int) {
			return bytes.LastIndex(data, end), len(end)
		}
		m.Spliter = SplitBetween(m.MatchStart, m.MatchEnd)
	}
	return m
}

//解析字节流
func SplitStream(reader io.Reader, spliter bufio.SplitFunc) (result [][]byte, err error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(spliter)
	defer func() {
		var ok bool
		if err, ok = recover().(error); !ok || err == nil {
			err = scanner.Err()
		}
	}()
	for scanner.Scan() {
		if chunk := scanner.Bytes(); chunk != nil {
			result = append(result, chunk)
		}
	}
	return
}
