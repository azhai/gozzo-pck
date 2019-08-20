package serialize

import (
	"bytes"
	"encoding/binary"

	"github.com/azhai/gozzo-utils/common"
)

// 单个字节
type Byte byte

func (n Byte) Encode(v interface{}) []byte {
	return []byte{v.(byte)}
}

func (n Byte) Decode(chunk []byte) interface{} {
	v := byte(0x00)
	if chunk != nil {
		v = chunk[len(chunk) - 1]
	}
	return v
}

// 字节数组
type Bytes []byte

func (s Bytes) Encode(v interface{}) []byte {
	return v.([]byte)
}

func (s Bytes) Decode(chunk []byte) interface{} {
	return chunk
}

// 字符串
type String string

func (s String) Encode(v interface{}) []byte {
	return []byte(v.(string))
}

func (s String) Decode(chunk []byte) interface{} {
	return string(chunk)
}

// BCD码
type HexStr string

func (s HexStr) Encode(v interface{}) []byte {
	return common.Hex2Bin(v.(string))
}

func (s HexStr) Decode(chunk []byte) interface{} {
	return common.Bin2Hex(chunk)
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
	buf := bytes.NewBuffer(nil)
	_ = binary.Write(buf, binary.BigEndian, v)
	chunk := make([]byte, n.Cap())
	if size, _ := buf.Read(chunk); size == n.Size {
		return chunk
	}
	return common.ResizeBytes(chunk, true, n.Size)
}

func (n *Unsigned) Decode(chunk []byte) interface{} {
	capSize := n.Cap()
	chunk = common.ResizeBytes(chunk, true, capSize)
	switch capSize {
	case 1:
		return chunk[0]
	case 2:
		return binary.BigEndian.Uint16(chunk)
	case 4:
		return binary.BigEndian.Uint32(chunk)
	default:
		return binary.BigEndian.Uint64(chunk)
	}
}

// 整数
type Integer struct {
	Negative bool
	*Unsigned
}

func (n *Integer) DecodeInt64(chunk []byte) (v int64) {
	u := n.Unsigned.Decode(chunk)
	switch u := u.(type) {
	default:
		return 0
	case uint8:
		v = int64(u)
	case uint16:
		v = int64(u)
	case uint32:
		v = int64(u)
	case uint64:
		v = int64(u)
	}
	if n.Negative {
		return 0 - v
	} else {
		return v
	}
}

func (n *Integer) Decode(chunk []byte) interface{} {
	v := n.DecodeInt64(chunk)
	switch n.Cap() {
	case 1:
		return int8(v)
	case 2:
		return int16(v)
	case 4:
		return int32(v)
	default:
		return v
	}
}
