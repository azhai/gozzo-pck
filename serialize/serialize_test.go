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
	remarks = []string{"成功/确认", "失败", "消息有误", "不支持", "报警处理确认"}
)

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

// JT/T808协议外层
type Proto808 struct {
	Head   byte    // 消息头
	Code   string  // 命令ID
	Props  uint16  // 属性，不加密不分包时即消息体长度
	Mobile string  // SIM卡号
	Seqno  uint16  // 流水号
	Rest   []byte  // 消息体（字段名称固定不可改），暂不解析，只获得原始字节
	Check  byte    // 检验码
	Tail   byte    // 消息尾
	*Object
}

func NewProto808() *Proto808 {
	p := &Proto808{
		Head: 0x7e, Tail: 0x7e,
		Object: NewObject(),
	}
	// 从前往后
	p.AddByteField("head", false)
	p.AddHexStrField("code", 2)
	p.AddUintField("props", 2)
	p.AddHexStrField("mobile", 6)
	p.AddUintField("seqno", 2)
	// 从后往前
	p.AddByteField("tail", true)
	p.AddByteField("check", true)
	// 剩余的会作为rest自动加上
	return p
}

func TestProto808(t *testing.T) {
	var seq = uint16(0)
	p := NewProto808()
	for _, msg := range []string{data808, reply808} {
		// 先还原后校验
		chunk := Unescape(common.Hex2Bin(msg))
		t.Log(common.Bin2Hex(chunk))
		assert.Equal(t, byte(0x00), BlockCheck(chunk))
		// 解析
		err := Unserialize(chunk, p)
		assert.NoError(t, err)
		assert.Equal(t, byte(0x7e), p.Head)
		assert.Equal(t, byte(0x7e), p.Tail)
		assert.Len(t, p.Rest, int(p.Props))
		assert.Equal(t, "082035085667", p.Mobile)
		seq = p.Seqno - seq
		t.Logf("%+v\n", p)

		if len(p.Rest) == 5 {
			_testBodyReply(t, p.Rest)
			assert.Equal(t, "8001", p.Code)
		} else {
			assert.Equal(t, uint16(7), p.Seqno)
			assert.Equal(t, "0200", p.Code)
		}
	}
	assert.Equal(t, uint16(1), seq)
}

// JT/T808协议，平台通用回复消息体
type BodyReply struct {
	Seqno  uint16  // 原消息的流水号
	Code   string  // 原消息的命令ID
	Status byte   // 结果 0:成功/确认；1:失败；2:消息有误；3:不支持；4:报警处理确认
	*Object
}

func NewBodyReply() *BodyReply {
	b := &BodyReply{Object: NewObject()}
	b.AddUintField("seqno", 2)
	b.AddHexStrField("code", 2)
	b.AddByteField("status", false)
	// opts := NewOptions(remarks)
	// b.AddEnumField("status", opts)
	return b
}

func _testBodyReply(t *testing.T, body []byte) {
	b := NewBodyReply()
	t.Log(common.Bin2Hex(body))
	// 解析
	err := Unserialize(body, b)
	assert.NoError(t, err)
	assert.Equal(t, uint16(7), b.Seqno)
	assert.Equal(t, "0200", b.Code)
	assert.Equal(t, byte(0x00), b.Status)
	assert.Equal(t, body, Serialize(b))
	t.Logf("%+v\n", b)
}
