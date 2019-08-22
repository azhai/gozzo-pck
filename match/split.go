package match

import (
	"bufio"
	"bytes"
	"io"
	"time"
)

// 判断匹配开头或结尾的方法
type MatchFunc func(data []byte) (int, int)

// 按前后标记拆包
type SplitMatcher struct {
	split bufio.SplitFunc
}

func NewSplitMatcher(split bufio.SplitFunc) *SplitMatcher {
	return &SplitMatcher{split: split}
}

// 解析字节流
func (m *SplitMatcher) SplitStream(rd io.Reader, outch chan<- []byte) error {
	scanner := bufio.NewScanner(rd)
	scanner.Split(m.split)
	for scanner.Scan() {
		if chunk := scanner.Bytes(); chunk != nil {
			outch <- chunk
		}
	}
	return scanner.Err()
}

// 解析二进制数据
func (m *SplitMatcher) SplitBuffer(input []byte) ([][]byte, error) {
	var output [][]byte
	rd := bytes.NewReader(input)
	outch := make(chan []byte)
	go func() {
		defer close(outch)
		for chunk := range outch {
			output = append(output, chunk)
		}
	}()
	err := m.SplitStream(rd, outch)
	return output, err
}

// 按前后标记分拆
type SplitCreator struct {
	StartToken []byte
	EndToken   []byte
}

func NewSplitCreator(start, end []byte) *SplitCreator {
	return &SplitCreator{StartToken: start, EndToken: end}
}

func (m *SplitCreator) GetSplit() bufio.SplitFunc {
	// 只有结尾标记
	if m.StartToken == nil {
		return SplitAfter(CreateMatchForward(m.EndToken))
	}
	matchStart := CreateMatchForward(m.StartToken)
	if bytes.Compare(m.StartToken, m.EndToken) == 0 {
		return SplitBoth(matchStart) // 前后标记相同
	} else {
		matchEnd := CreateMatchBackward(m.EndToken)
		return SplitBetween(matchStart, matchEnd) // 前后标记不同
	}
}

func CreateMatchForward(token []byte) MatchFunc {
	return func(data []byte) (int, int) {
		return bytes.Index(data, token), len(token)
	}
}

func CreateMatchBackward(token []byte) MatchFunc {
	return func(data []byte) (int, int) {
		return bytes.LastIndex(data, token), len(token)
	}
}

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

// 根据结尾分割，但是从前向后搜索
func SplitAfter(match MatchFunc) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (int, []byte, error) {
		i, m := match(data)
		if i >= 0 {
			return i + m, data[:i+m], nil
		}
		var err error
		if atEOF {
			err = bufio.ErrFinalToken
		}
		return 0, nil, err
	}
}

// 根据相同的开头和结尾分割
func SplitBoth(match MatchFunc) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (int, []byte, error) {
		advance, i, m := MatchTwice(match, data, atEOF)
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

// 按定长（可选定时）分拆
type FixedSplitCreator struct {
	ByteSize int
	Interval time.Duration
}

func NewFixedSplitCreator(size int, msec int64) *FixedSplitCreator {
	return &FixedSplitCreator{
		ByteSize: size,
		Interval: time.Duration(msec) * time.Millisecond,
	}
}

func (m *FixedSplitCreator) GetSplit() bufio.SplitFunc {
	return func(data []byte, atEOF bool) (int, []byte, error) {
		timer := time.NewTimer(m.Interval)
		for {
			select {
			default:
				if len(data) >= m.ByteSize {
					timer.Stop()
					return m.ByteSize, data[:m.ByteSize], nil
				}
				if atEOF { // ErrFinalToken，最后一段数据长度不够也可以结束
					return len(data), data, bufio.ErrFinalToken
				}
				return 0, nil, nil
			case <- timer.C:
				timer.Reset(m.Interval)
				return len(data), data, nil
			}
		}
	}
}
