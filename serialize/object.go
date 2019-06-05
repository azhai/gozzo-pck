package serialize

import (
	"reflect"
	"strings"

	"github.com/azhai/gozzo-pck/match"
)

type ISerializer interface {
	Serialize(obj interface{}) []byte
	Unserialize(chunk []byte, objref interface{}) error
}

// 对象
type Object struct {
	children map[string]IEncoder
	Matcher  *match.FieldMatcher
}

func NewObject() *Object {
	t := &Object{
		children: make(map[string]IEncoder),
		Matcher:  match.NewFieldMatcher(),
	}
	t.children["rest"] = new(Bytes)
	return t
}

func (t *Object) Serialize(obj interface{}) []byte {
	data := make(map[string][]byte)
	rv := reflect.Indirect(reflect.ValueOf(obj))
	var val reflect.Value
	for name, child := range t.children {
		val = rv.FieldByName(strings.Title(name))
		if val.IsValid() { // 存在的字段
			data[name] = child.Encode(val.Interface())
		}
	}
	return t.Matcher.Build(data)
}

func (t *Object) Unserialize(chunk []byte, objref interface{}) error {
	data := t.Matcher.Match(chunk, true)
	rv := reflect.Indirect(reflect.Indirect(reflect.ValueOf(objref)))
	var val reflect.Value
	for name, child := range t.children {
		if bin, ok := data[name]; ok && bin != nil {
			val = rv.FieldByName(strings.Title(name))
			if val.IsValid() && val.CanSet() { // 存在的字段
				val.Set(reflect.ValueOf(child.Decode(bin)))
			}
		}
	}
	return nil
}

func (t *Object) Encode(v interface{}) []byte {
	return t.Serialize(v)
}

func (t *Object) Decode(chunk []byte) interface{} {
	if err := t.Unserialize(chunk, &t); err == nil {
		return t
	}
	return nil
}

func (t *Object) GetChild(name string) (IEncoder, bool) {
	child, ok := t.children[name]
	return child, ok
}

func (t *Object) AddChild(child IEncoder, name string, size int, rev bool) *match.Field {
	t.children[name] = child
	if rev && size > 0 {
		size = 0 - size
	}
	field := match.NewField(size, false)
	return t.Matcher.AddField(name, field)
}

func (t *Object) AddSpanField(size int) *match.Field {
	return t.Matcher.AddField("", match.NewField(size, false))
}

func (t *Object) AddByteField(name string, rev bool) *match.Field {
	return t.AddChild(new(Byte), name, 1, rev)
}

func (t *Object) AddBytesField(name string, size int) *match.Field {
	return t.AddChild(new(Bytes), name, size, false)
}

func (t *Object) AddStringField(name string, size int) *match.Field {
	return t.AddChild(new(String), name, size, false)
}

func (t *Object) AddHexStrField(name string, size int) *match.Field {
	return t.AddChild(new(HexStr), name, size, false)
}

func (t *Object) AddUint64Field(name string, rev bool) *match.Field {
	return t.AddChild(new(Uint64), name, 8, rev)
}

func (t *Object) AddUint32Field(name string, rev bool) *match.Field {
	return t.AddChild(new(Uint32), name, 4, rev)
}

func (t *Object) AddUint24Field(name string, rev bool) *match.Field {
	return t.AddChild(new(Uint24), name, 3, rev)
}

func (t *Object) AddUint16Field(name string, rev bool) *match.Field {
	return t.AddChild(new(Uint16), name, 2, rev)
}

func (t *Object) AddEnumField(name string, rev bool, opts *Options) *match.Field {
	return t.AddChild(NewEnum(opts), name, 1, rev)
}
