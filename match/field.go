package match

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
func (field *Field) GetRange(offset int) (int, int) {
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
	return start, stop
}

//分段匹配
type FieldMatcher struct {
	Sequence []*Field //开头已定义段
	Reverse  []*Field //结尾已定义段
	Rest     *Field   //未识别部分，可作为payload创建新包
}

func NewFieldMatcher() *FieldMatcher {
	rest := NewField(0, false) //初始时，全部字节未识别
	return &FieldMatcher{Rest: rest}
}

func (m *FieldMatcher) AddFixeds(sizes []int) {
	if len(sizes) > 0 {
		for _, size := range sizes {
			m.AddField(NewField(size, false))
		}
	}
}

//添加开头的段定义
func (m *FieldMatcher) AddField(field *Field) *Field {
	field.Start += m.Rest.Start
	if field.Size > 0 {
		field.Stop = field.Start + field.Size
		if !field.Optional { //增加固定字段时，缩减未知部分的范围
			m.Rest.Start = field.Stop
		}
	}
	m.Sequence = append(m.Sequence, field)
	return field
}

//添加结尾的段定义
func (m *FieldMatcher) AddRevField(field *Field) *Field {
	field.Start += m.Rest.Stop
	if field.Size > 0 {
		field.Stop = field.Start + field.Size
		if !field.Optional { //缩减未知部分的范围
			m.Rest.Stop = field.Start
		}
	}
	m.Reverse = append(m.Reverse, field)
	return field
}

func (m *FieldMatcher) GetLeastSize() (int, int) {
	var num, least = 0, 0
	for _, f := range m.Sequence {
		num++
		least += f.Size
	}
	if len(m.Reverse) > 0 {
		for _, f := range m.Reverse {
			num++
			least += f.Size
		}
	}
	return num, least
}

// 直接定义并读取开头几个固定段，类似Erlang中的位匹配
func MatchHead(chunk []byte, squence []*Field) ([][]byte, int) {
	var (
		result      [][]byte
		start, stop int
	)
	size := len(chunk)
	for _, field := range squence {
		start, stop = field.GetRange(0) // 一定不为负
		if start >= 0 && stop <= size {
			result = append(result, chunk[start:stop])
		} else {
			result = append(result, nil)
		}
	}
	return result, stop
}
