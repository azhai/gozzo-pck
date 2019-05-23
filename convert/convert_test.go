package convert

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func Unescape(data []byte) []byte {
	data = bytes.Replace(data, []byte{0x7d, 0x02}, []byte{0x7e}, -1)
	return bytes.Replace(data, []byte{0x7d, 0x01}, []byte{0x7d}, -1)
}

func GetT808Body(data string) (chunk []byte) {
	chunk, _ = hex.DecodeString(data)
	chunk = Unescape(chunk[1 : len(chunk)-1])
	return
}

// 测试JT/T808协议
func TestObject(t *testing.T) {
	data := "7E01020006014530399195003F717361757468597E"
	p := NewObject()
	p.AddHexStrField("code", 2)
	p.AddUint16Field("props")
	p.AddHexStrField("mobile", 6)
	p.AddUint16Field("msgno")
	chunk := GetT808Body(data)
	t.Log(hex.EncodeToString(chunk))
	p.Decode(chunk)
	t.Log("code", p.GetHexStr("code"))
	t.Log("mobile", p.GetHexStr("mobile"))
	t.Log("msgno", p.GetUint16("msgno"))
	t.Log("props", p.GetUint16("props"))
}
