package match

import (
	"fmt"

	"github.com/azhai/gozzo-utils/common"
)

//段，若干个byte组成
type Field struct {
	Size     int  //长度>=0
	Optional bool //可选或不定长度（非固定）
	Start    int  //开始位置（包含），可能为负
	Stop     int  //结束位置（不包含），可能为负
}

func NewField(size int, optional bool) *Field {
	start := 0
	if size < 0 {
		start = size
		size = 0 - size //取绝对值
	}
	return &Field{
		Size: size, Optional: optional,
		Start: start, Stop: 0,
	}
}

//找出段的正向起止位置，offset为修正值，只对同符号数据起作用
func (field *Field) GetRange(offset, size int) (int, int) {
	var (
		start = field.Start
		stop  = field.Stop
	)
	if start*offset > 0 { //皆为正或皆为负，加上修正值
		start += offset
	}
	if stop*offset >= 0 { //同符号时修正
		stop += offset
	}
	if stop <= 0 {
		stop += size
	}
	return start, stop
}

//分段匹配
type FieldMatcher struct {
	rest     *Field //未识别部分，可作为payload创建新包
	fields   map[string]*Field
	Sequence []string //开头已定义段
	Reverse  []string //结尾已定义段
}

func NewFieldMatcher() *FieldMatcher {
	return &FieldMatcher{
		rest:   NewField(0, false), //初始时，全部字节未识别
		fields: make(map[string]*Field),
	}
}

func (m *FieldMatcher) AddFixeds(sizes []int) {
	if len(sizes) > 0 {
		for _, size := range sizes {
			m.AddField("", NewField(size, false))
		}
	}
}

//添加开头的段定义
func (m *FieldMatcher) AddField(name string, field *Field) *Field {
	if name == "" || name == "rest" {
		name = fmt.Sprintf("+%d", len(m.Sequence))
	}
	field.Start += m.rest.Start
	if field.Size > 0 {
		field.Stop = field.Start + field.Size
		if !field.Optional { //增加固定字段时，缩减未知部分的范围
			m.rest.Start = field.Stop
		}
	}
	m.Sequence = append(m.Sequence, name)
	m.fields[name] = field
	return field
}

//添加结尾的段定义
func (m *FieldMatcher) AddRevField(name string, field *Field) *Field {
	if name == "" || name == "rest" {
		name = fmt.Sprintf("-%d", len(m.Sequence)+1)
	}
	field.Start += m.rest.Stop
	if field.Size > 0 {
		field.Stop = field.Start + field.Size
		if !field.Optional { //缩减未知部分的范围
			m.rest.Stop = field.Start
		}
	}
	m.Reverse = append(m.Reverse, name)
	m.fields[name] = field
	return field
}

func (m *FieldMatcher) GetLeastSize() (int, int) {
	var least = 0
	for _, f := range m.fields {
		least += f.Size
	}
	return len(m.fields), least
}

// 直接定义并读取开头几个固定段，类似Erlang中的位匹配
func (m *FieldMatcher) Match(chunk []byte, withRest bool) map[string][]byte {
	var data = make(map[string][]byte)
	start, stop, size := 0, 0, len(chunk)
	for name, field := range m.fields {
		start, stop = field.GetRange(0, size) // 一定不为负
		if start >= 0 && stop <= size {
			data[name] = chunk[start:stop]
		} else {
			data[name] = nil
		}
	}
	if withRest {
		data["rest"] = nil
		start, stop = m.rest.GetRange(0, size)
		if start >= 0 && stop <= size {
			data["rest"] = chunk[start:stop]
		}
	}
	return data
}

func (m *FieldMatcher) Build(data map[string][]byte) []byte {
	var (
		field        *Field
		ok           bool
		chunk, value []byte
	)
	names := append(m.Sequence, "rest")
	names = append(names, m.Reverse...)
	for _, name := range names {
		if field, ok = m.fields[name]; !ok {
			continue
		}
		if value, ok = data[name]; !ok {
			value = nil
		}
		if field.Size > 0 {
			value = common.ResizeBytes(value, true, field.Size)
		}
		chunk = append(chunk, value...)
	}
	return chunk
}
