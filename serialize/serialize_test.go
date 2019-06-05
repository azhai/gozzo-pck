package serialize

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
	obj := NewObject()
	obj.AddHexStrField("code", 2)
	obj.AddUint16Field("props", false)
	obj.AddHexStrField("mobile", 6)
	obj.AddUint16Field("msgno", false)
	obj.AddByteField("check", true)
	chunk := GetT808Body(data)
	t.Log(hex.EncodeToString(chunk))
	obj.Unserialize(chunk, &obj)
	/*t.Log("code", obj.Table["code"].(string))
	t.Log("mobile", obj.Table["mobile"].(string))
	t.Log("msgno", obj.Table["msgno"].(uint16))
	t.Log("props", obj.Table["props"].(uint16))*/
}
