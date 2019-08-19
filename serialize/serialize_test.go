package serialize

import (
	"bytes"
	"testing"

	"github.com/azhai/gozzo-utils/common"
	"github.com/stretchr/testify/assert"
)

var (
	data808 = "7e0200005b08203508566700070000" +
		"0000000c000301c905250741c82b00030000011c1907200402470104000000fa" +
		"2a02000030011831010c57080000000000000000fb0100fc02000afd0209b5fe" +
		"143839383630343132313031393931303233313532ff020001917e"
	reply808 = "7e8001000508203508566700080007020000ad7e"
)

// JT/T808协议解析
type Proto808 struct {
	Head   byte
	Code   string
	Props  uint16
	Mobile string
	Seqno  uint16
	Check  byte
	Tail   byte
	Rest   []byte
	*Object
}

func NewProto808() *Proto808 {
	obj := &Proto808{
		Head: 0x7e, Tail: 0x7e,
		Object: NewObject(),
	}
	// 从前往后
	obj.AddByteField("head", false)
	obj.AddHexStrField("code", 2)
	obj.AddUintField("props", 2)
	obj.AddHexStrField("mobile", 6)
	obj.AddUintField("seqno", 2)
	// 从后往前
	obj.AddByteField("tail", true)
	obj.AddByteField("check", true)
	// 剩余的会作为rest自动加上
	return obj
}

func Escape(data []byte) []byte {
	data = bytes.Replace(data, []byte{0x7d}, []byte{0x7d, 0x01}, -1)
	return bytes.Replace(data, []byte{0x7e}, []byte{0x7d, 0x02}, -1)
}

func Unescape(data []byte) []byte {
	data = bytes.Replace(data, []byte{0x7d, 0x02}, []byte{0x7e}, -1)
	return bytes.Replace(data, []byte{0x7d, 0x01}, []byte{0x7d}, -1)
}

//异或校验
func BlockCheck(block []byte) byte {
	result := byte(0x00)
	for _, bin := range block {
		result ^= bin
	}
	return result
}

func TestProto808(t *testing.T) {
	var seq = uint16(0)
	obj := NewProto808()
	for _, msg := range []string{data808, reply808} {
		// 先还原后校验
		chunk := Unescape(common.Hex2Bin(msg))
		t.Log(common.Bin2Hex(chunk))
		assert.Equal(t, byte(0x00), BlockCheck(chunk))
		// 解析
		_, err := obj.Decode(chunk)
		assert.NoError(t, err)
		assert.Equal(t, byte(0x7e), obj.Head)
		assert.Equal(t, byte(0x7e), obj.Tail)
		assert.Len(t, obj.Rest, int(obj.Props))
		assert.Equal(t, "082035085667", obj.Mobile)
		seq = obj.Seqno - seq
		t.Logf("%+v\n", obj)
	}
	assert.Equal(t, uint16(1), seq)
}
