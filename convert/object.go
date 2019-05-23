package convert

import (
	"sort"

	"github.com/azhai/gozzo-pck/match"
)

// 对象
type Object struct {
	current  []string
	children map[string]IConvert
	Matcher  *match.FieldMatcher
}

func NewObject() *Object {
	p := &Object{
		Matcher: match.NewFieldMatcher(),
		children: make(map[string]IConvert),
	}
	p.children["rest"] = new(Bytes)
	return p
}

func (p *Object) GetConvert(name string) (IConvert, bool) {
	child, ok := p.children[name]
	return child, ok
}

func (p *Object) Encode() []byte {
	var data = make(map[string][]byte)
	for name, child := range p.children {
		data[name] = child.Encode()
	}
	return p.Matcher.Build(data)
}

func (p *Object) Decode(chunk []byte) (err error) {
	p.current = make([]string, 0)
	data := p.Matcher.Match(chunk, true)
	for name, value := range data {
		if child, ok := p.children[name]; ok {
			p.current = append(p.current, name)
			err = child.Decode(value)
		}
	}
	sort.Strings(p.current)
	return
}

func (p *Object) GetData() interface{} {
	return p.GetTable()
}

func (p *Object) SetData(data interface{}) {
	if table, ok := data.(map[string]interface{}); ok {
		p.SetTable(table)
	}
}

func (p *Object) GetTable() map[string]interface{} {
	result := make(map[string]interface{})
	for _, name := range p.current {
		result[name] = p.children[name].GetData()
	}
	return result
}

func (p *Object) SetTable(table map[string]interface{}) {
	p.current = make([]string, 0)
	for name, conv := range p.children {
		if data, ok := table[name]; ok {
			p.current = append(p.current, name)
			conv.SetData(data)
		}
	}
	sort.Strings(p.current)
}

func (p *Object) AddSpanField(size int) *match.Field {
	return p.Matcher.AddField("", match.NewField(size, false))
}

func (p *Object) AddBytesField(name string, size int) *match.Field {
	if _, ok := p.children[name]; ok {
		return nil
	}
	p.children[name] = new(Bytes)
	return p.Matcher.AddField(name, match.NewField(size, false))
}

func (p *Object) GetBytes(name string) (value []byte) {
	if conv, ok := p.children[name]; ok {
		if s, succ := conv.(*Bytes); succ {
			value = s.Data
		}
	}
	return
}

func (p *Object) AddStringField(name string, size int) *match.Field {
	if _, ok := p.children[name]; ok {
		return nil
	}
	p.children[name] = new(String)
	return p.Matcher.AddField(name, match.NewField(size, false))
}

func (p *Object) GetString(name string) (value string) {
	if conv, ok := p.children[name]; ok {
		if s, succ := conv.(*String); succ {
			value = s.Data
		}
	}
	return
}

func (p *Object) AddHexStrField(name string, size int) *match.Field {
	if _, ok := p.children[name]; ok {
		return nil
	}
	p.children[name] = new(HexStr)
	return p.Matcher.AddField(name, match.NewField(size, false))
}

func (p *Object) GetHexStr(name string) (value string) {
	if conv, ok := p.children[name]; ok {
		if s, succ := conv.(*HexStr); succ {
			value = s.Data
		}
	}
	return
}

func (p *Object) AddUint64Field(name string) *match.Field {
	if _, ok := p.children[name]; ok {
		return nil
	}
	p.children[name] = new(Uint64)
	return p.Matcher.AddField(name, match.NewField(8, false))
}

func (p *Object) GetUint64(name string) (value uint64) {
	if conv, ok := p.children[name]; ok {
		if s, succ := conv.(*Uint64); succ {
			value = s.Data
		}
	}
	return
}

func (p *Object) AddUint32Field(name string) *match.Field {
	if _, ok := p.children[name]; ok {
		return nil
	}
	p.children[name] = new(Uint32)
	return p.Matcher.AddField(name, match.NewField(4, false))
}

func (p *Object) GetUint32(name string) (value uint32) {
	if conv, ok := p.children[name]; ok {
		if s, succ := conv.(*Uint32); succ {
			value = s.Data
		}
	}
	return
}

func (p *Object) AddUint24Field(name string) *match.Field {
	if _, ok := p.children[name]; ok {
		return nil
	}
	p.children[name] = new(Uint24)
	return p.Matcher.AddField(name, match.NewField(3, false))
}

func (p *Object) GetUint24(name string) (value uint32) {
	if conv, ok := p.children[name]; ok {
		if s, succ := conv.(*Uint24); succ {
			value = s.Uint32.Data
		}
	}
	return
}

func (p *Object) AddUint16Field(name string) *match.Field {
	if _, ok := p.children[name]; ok {
		return nil
	}
	p.children[name] = new(Uint16)
	return p.Matcher.AddField(name, match.NewField(2, false))
}

func (p *Object) GetUint16(name string) (value uint16) {
	if conv, ok := p.children[name]; ok {
		if s, succ := conv.(*Uint16); succ {
			value = s.Data
		}
	}
	return
}

func (p *Object) AddEnumField(name string, opts *Options) *match.Field {
	if _, ok := p.children[name]; ok {
		return nil
	}
	p.children[name] = &Enum{Options:opts}
	return p.Matcher.AddField(name, match.NewField(1, false))
}

func (p *Object) GetEnum(name string) *Enum {
	if conv, ok := p.children[name]; ok {
		if s, succ := conv.(*Enum); succ {
			return s
		}
	}
	return nil
}

func (p *Object) GetEnumByte(name string) (value byte) {
	if s := p.GetEnum(name); s != nil {
		value = s.Data
	}
	return
}

func (p *Object) GetEnumString(name string) (remark string) {
	if s := p.GetEnum(name); s != nil {
		remark = s.ToString()
	}
	return
}
