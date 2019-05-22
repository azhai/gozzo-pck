package convert

import (
	"sort"
	"strings"
)

// 枚举选项
type Options struct {
	values  []byte
	remarks []string
}

func NewOptions(remarks []string) *Options {
	opts := &Options{remarks: remarks}
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
	if i := t.Search(v); i >= 0 && i < t.Size() {
		if t.values[i] == v {
			return i
		}
	}
	return -1
}

func (t *Options) ByRemark(r string) int {
	for i, mark := range t.remarks {
		if strings.Compare(mark, r) == 0 {
			return i
		}
	}
	return -1
}

// 枚举类型
type Enum struct {
	Data byte
	Index int
	*Options
}

func (n *Enum) Encode() []byte {
	return []byte{n.Data}
}

func (n *Enum) Decode(chunk []byte) error {
	if chunk != nil {
		n.Data = chunk[0]
	}
	return nil
}

func (n *Enum) GetData() interface{} {
	return n.Data
}

func (n *Enum) SetData(data interface{}) {
	var val = data.(byte)
	if i := n.ByValue(val); i >= 0 {
		n.Index = i
		n.Data = val
	}
}

func (n *Enum) SetIndex(i int) {
	if i >= 0 && i < n.Size() {
		n.Index = i
		n.Data, _ = n.Item(i)
	}
}

func (n *Enum) SetString(mark string) {
	if i := n.ByRemark(mark); i >= 0 {
		n.Index = i
		n.Data, _ = n.Item(i)
	}
}

func (n *Enum) ToString() string {
	val, mark := n.Item(n.Index)
	if n.Data == val {
		return mark
	}
	return ""
}
