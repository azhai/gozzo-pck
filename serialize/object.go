package serialize

import (
	"reflect"
	"strings"

	"github.com/azhai/gozzo-pck/match"
)

type IEncoder interface {
	Encode(v interface{}) []byte
	Decode(chunk []byte) interface{}
}

type ISerializer interface {
	GetMatcher() *match.FieldMatcher
	GetNames() map[string]string
	GetChild(name string) (IEncoder, bool)
}

func Serialize(s ISerializer) []byte {
	data := make(map[string][]byte)
	rv := reflect.Indirect(reflect.ValueOf(s))
	for name, prop := range s.GetNames() {
		child, ok := s.GetChild(name)
		rf := rv.FieldByName(prop)
		if ok && rf.IsValid() { // 存在的字段
			data[name] = child.Encode(rf.Interface())
		}
	}
	return s.GetMatcher().Build(data)
}

func Unserialize(chunk []byte, s ISerializer) (err error) {
	data := s.GetMatcher().Match(chunk, true)
	rv := reflect.Indirect(reflect.ValueOf(s))
	var val interface{}
	for name, prop := range s.GetNames() {
		child, ok := s.GetChild(name)
		rf := rv.FieldByName(prop)
		if !ok || !rf.IsValid() || !rf.CanSet() {
			continue
		}
		if bin, ok := data[name]; ok {
			val = child.Decode(bin)
		} else {
			val = child.Decode(nil)
		}
		rf.Set(reflect.ValueOf(val))
	}
	return
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

func (t *Object) GetMatcher() *match.FieldMatcher {
	return t.Matcher
}

func (t *Object) GetNames() map[string]string {
	names := make(map[string]string)
	for name := range t.children {
		names[name] = strings.Title(name)
	}
	return names
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
	return t.AddFixedChild(name, NewUnsigned(size), size, false)
}

func (t *Object) AddEnumField(name string, opts *Options) (*match.Field, *Enum) {
	m := NewEnum(opts)
	f := t.AddFixedChild(name, m, 1, false)
	return f, m
}

func (t *Object) AddTwoDimField(name string, size int, x, y int64) (*match.Field, *TwoDim) {
	td := NewTwoDimXY(size, x, y)
	f := t.AddFixedChild(name, td, td.Size*2, false)
	return f, td
}

func (t *Object) AddTimeStampField(name string) (*match.Field, *TimeStamp) {
	ts := NewTimeStamp()
	f := t.AddFixedChild(name, ts, ts.Size, false)
	return f, ts
}

func (t *Object) AddDateField(name string) (*match.Field, *Date) {
	d := new(Date)
	f := t.AddFixedChild(name, d, 4, false)
	return f, d
}
