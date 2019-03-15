package convert

import "github.com/azhai/gozzo-pck/match"

// 对象
type Object struct {
	names    []string
	children map[string]IConvert
	Matcher  *match.FieldMatcher
}

func (p *Object) Init() *Object {
	p.names = make([]string, 0)
	p.children = make(map[string]IConvert)
	p.Matcher = match.NewFieldMatcher()
	return p
}

func (p *Object) GetConvert(name string) (IConvert, bool) {
	conv, ok := p.children[name]
	return conv, ok
}

func (p *Object) Encode() []byte {
	var result []byte
	for _, name := range p.names {
		if name == "" {
			continue
		}
		if value := p.children[name].Encode(); value != nil {
			result = append(result, value...)
		}
	}
	return result
}

func (p *Object) Decode(chunk []byte) error {
	var (
		err  error
		name string
		size = len(chunk)
	)
	for i, field := range p.Matcher.Sequence {
		if name = p.names[i]; name == "" {
			continue
		}
		start, stop := field.GetRange(0)
		if start >= 0 && stop <= size {
			err = p.children[name].Decode(chunk[start:stop])
		}
	}
	return err
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
	for name, conv := range p.children {
		result[name] = conv.GetData()
	}
	return result
}

func (p *Object) SetTable(table map[string]interface{}) {
	for name, conv := range p.children {
		if data, ok := table[name]; ok {
			conv.SetData(data)
		}
	}
}

func (p *Object) AddSpanField(size int) *match.Field {
	p.names = append(p.names, "")
	return p.Matcher.AddField(match.NewField(size, false))
}

func (p *Object) AddBytesField(name string, size int) *match.Field {
	if _, ok := p.children[name]; ok {
		return nil
	}
	p.names = append(p.names, name)
	p.children[name] = new(Bytes)
	return p.Matcher.AddField(match.NewField(size, false))
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
	p.names = append(p.names, name)
	p.children[name] = new(String)
	return p.Matcher.AddField(match.NewField(size, false))
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
	p.names = append(p.names, name)
	p.children[name] = new(HexStr)
	return p.Matcher.AddField(match.NewField(size, false))
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
	p.names = append(p.names, name)
	p.children[name] = new(Uint64)
	return p.Matcher.AddField(match.NewField(8, false))
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
	p.names = append(p.names, name)
	p.children[name] = new(Uint32)
	return p.Matcher.AddField(match.NewField(4, false))
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
	p.names = append(p.names, name)
	p.children[name] = new(Uint24)
	return p.Matcher.AddField(match.NewField(3, false))
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
	p.names = append(p.names, name)
	p.children[name] = new(Uint16)
	return p.Matcher.AddField(match.NewField(2, false))
}

func (p *Object) GetUint16(name string) (value uint16) {
	if conv, ok := p.children[name]; ok {
		if s, succ := conv.(*Uint16); succ {
			value = s.Data
		}
	}
	return
}
