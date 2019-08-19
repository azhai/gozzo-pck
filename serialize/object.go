package serialize

import (
	"reflect"
	"strings"

	"github.com/azhai/gozzo-pck/match"
)

type ISerializer interface {
	AddChild(name string, child IEncoder, field *match.Field)
	GetChild(name string) (IEncoder, bool)
	IEncoder
}

// 对象
type Object struct {
	Translate func(s string) string
	children  map[string]IEncoder
	Matcher   *match.FieldMatcher
}

func NewObject() *Object {
	t := &Object{
		Translate: strings.Title,
		children: make(map[string]IEncoder),
		Matcher:  match.NewFieldMatcher(),
	}
	t.children["rest"] = new(Bytes)
	return t
}

func (t Object) Encode(v interface{}) []byte {
	data := make(map[string][]byte)
	rv := reflect.Indirect(reflect.ValueOf(v))
	var val reflect.Value
	for name, child := range t.children {
		val = rv.FieldByName(t.Translate(name))
		if val.IsValid() { // 存在的字段
			data[name] = child.Encode(val.Interface())
		}
	}
	return t.Matcher.Build(data)
}

func (t *Object) Decode(chunk []byte) (interface{}, error) {
	var (
		err error
		v   interface{}
	)
	data := t.Matcher.Match(chunk, true)
	rv := reflect.Indirect(reflect.ValueOf(t))
	rv = reflect.Indirect(rv)
	for name, child := range t.children {
		if bin, ok := data[name]; ok && bin != nil {
			val := rv.FieldByName(strings.Title(name))
			if !val.IsValid() || !val.CanSet() { // 存在的字段
				continue
			}
			if v, err = child.Decode(bin); err == nil {
				val.Set(reflect.ValueOf(v))
			}
		}
	}
	return t, err
}

func (t *Object) GetChild(name string) (IEncoder, bool) {
	child, ok := t.children[name]
	return child, ok
}

func (t *Object) AddChild(name string, child IEncoder, field *match.Field) {
	if name != "" {
		t.children[name] = child
	}
	t.Matcher.AddField(name, field)
}

func (t *Object) AddFixedChild(name string, child IEncoder, size int, rev bool) *match.Field {
	if rev && size > 0 {
		size = 0 - size
	}
	field := match.NewField(size, false)
	t.AddChild(name, child, field)
	return field
}

func (t *Object) AddSpanField(size int, rev bool) *match.Field {
	return t.AddFixedChild("", nil, size, rev)
}

func (t *Object) AddByteField(name string, rev bool) *match.Field {
	return t.AddFixedChild(name, new(Byte), 1, rev)
}

func (t *Object) AddBytesField(name string, size int, rev bool) *match.Field {
	return t.AddFixedChild(name, new(Bytes), size, rev)
}

func (t *Object) AddStringField(name string, size int) *match.Field {
	return t.AddFixedChild(name, new(String), size, false)
}

func (t *Object) AddHexStrField(name string, size int) *match.Field {
	return t.AddFixedChild(name, new(HexStr), size, false)
}

func (t *Object) AddUintField(name string, size int) *match.Field {
	return t.AddFixedChild(name, &Unsigned{Size:size}, size, false)
}

func (t *Object) AddEnumField(name string, opts *Options) *match.Field {
	return t.AddFixedChild(name, NewEnum(opts), 1, false)
}
