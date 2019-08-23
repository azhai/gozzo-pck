package serialize

import (
	"fmt"
	"sort"
	"strings"
)

// 枚举选项
type Options struct {
	values     []byte
	remarks    []string
	isZeroBase bool // index和value的值相等（类型不同）
}

func NewOptions(remarks []string) *Options {
	opts := &Options{remarks: remarks, isZeroBase: true}
	for i := 0; i < len(opts.remarks); i++ {
		opts.values = append(opts.values, byte(i))
	}
	return opts
}

func NewMapOptions(options map[int]string) *Options {
	var values []int
	for val := range options {
		values = append(values, val)
	}
	sort.Ints(values)
	opts := &Options{}
	for _, val := range values {
		opts.values = append(opts.values, byte(val))
		opts.remarks = append(opts.remarks, options[val])
	}
	return opts
}

func (t *Options) Size() int {
	return len(t.values)
}

func (t *Options) Item(i int) (byte, string) {
	if i >= 0 && i < t.Size() {
		return t.values[i], t.remarks[i]
	}
	return 0x00, ""
}

func (t *Options) Search(v byte) int {
	return sort.Search(t.Size(), func(n int) bool {
		return t.values[n] >= v
	})
}

func (t *Options) ByValue(v byte) int {
	if t.isZeroBase {
		return int(v)
	}
	if i := t.Search(v); i >= 0 && i < t.Size() {
		if t.values[i] == v {
			return i
		}
	}
	return -1
}

func (t *Options) ByRemark(r string, caseSensitive bool) int {
	if caseSensitive {
		for i, mark := range t.remarks {
			if strings.Compare(mark, r) == 0 {
				return i
			}
		}
	} else { // 不区分大小写
		for i, mark := range t.remarks {
			if strings.EqualFold(mark, r) {
				return i
			}
		}
	}
	return -1
}

// 枚举类型
type Enum struct {
	Opts *Options
	Byte // opts的索引，范围0~255
}

func NewEnum(opts *Options) *Enum {
	return &Enum{Opts: opts, Byte: Byte(0)}
}

func (m Enum) GetIndex() int {
	return int(m.Byte)
}

func (m Enum) GetItem() (byte, string) {
	return m.Opts.Item(m.GetIndex())
}

func (m Enum) ToString() string {
	_, mark := m.GetItem()
	return mark
}

func (m *Enum) SetIndex(i int) error {
	if i >= 0 && i < m.Opts.Size() {
		m.Byte = Byte(i)
		return nil
	}
	return fmt.Errorf("Can not found the options of %d", i)
}

func (m *Enum) Decode(chunk []byte) interface{} {
	b := m.Byte.Decode(chunk).(byte)
	m.Byte = Byte(b)
	return b
}
