package serialize

import (
	"bytes"
	"encoding/binary"

	"github.com/azhai/gozzo-utils/common"
)

type IEncoder interface {
	Encode(v interface{}) []byte
	Decode(chunk []byte) (interface{}, error)
}

// 单个字节
type Byte byte

func (n Byte) Encode(v interface{}) []byte {
	return []byte{v.(byte)}
}

func (n Byte) Decode(chunk []byte) (interface{}, error) {
	v := byte(0x00)
	if chunk != nil {
		v = chunk[len(chunk) - 1]
	}
	return v, nil
}

// 字节数组
type Bytes []byte

func (s Bytes) Encode(v interface{}) []byte {
	return v.([]byte)
}

func (s Bytes) Decode(chunk []byte) (interface{}, error) {
	return chunk, nil
}

// 字符串
type String string

func (s String) Encode(v interface{}) []byte {
	return []byte(v.(string))
}

func (s String) Decode(chunk []byte) (interface{}, error) {
	return string(chunk), nil
}

// BCD码
type HexStr string

func (s HexStr) Encode(v interface{}) []byte {
	return common.Hex2Bin(v.(string))
}

func (s HexStr) Decode(chunk []byte) (interface{}, error) {
	v := common.Bin2Hex(chunk)
	return v, nil
}

// 无符号整数
type Unsigned struct {
	Size int
}

func (n Unsigned) Cap() int {
	// 修正错误的长度
	if n.Size < 1 {
		n.Size = 1
	} else if n.Size > 8 {
		n.Size = 8
	}
	// 对应合适的uint
	if n.Size == 3 {
		return 4
	} else if n.Size >= 5 && n.Size <= 7 {
		return 8
	} else {
		return n.Size
	}
}

func (n Unsigned) Encode(v interface{}) []byte {
	chunk := make([]byte, n.Cap())
	buf := bytes.NewBuffer(chunk)
	_ = binary.Write(buf, binary.BigEndian, v)
	return common.ResizeBytes(chunk, true, n.Size)
}

func (n Unsigned) Decode(chunk []byte) (interface{}, error) {
	chunk = common.ResizeBytes(chunk, true, n.Cap())
	buf := bytes.NewReader(chunk)
	err := binary.Read(buf, binary.BigEndian, &n)
	return n, err
}
