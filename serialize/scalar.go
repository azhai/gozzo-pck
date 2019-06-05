package serialize

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/azhai/gozzo-utils/common"
)

type IEncoder interface {
	Encode(v interface{}) []byte
	Decode(chunk []byte) interface{}
}

// 单个字节
type Byte byte

func (n Byte) Encode(v interface{}) []byte {
	return []byte{v.(byte)}
}

func (n Byte) Decode(chunk []byte) interface{} {
	if chunk != nil {
		return chunk[len(chunk) - 1]
	}
	return 0x00
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
	if chunk, err := hex.DecodeString(v.(string)); err == nil {
		return chunk
	}
	return nil
}

func (s HexStr) Decode(chunk []byte) interface{} {
	return hex.EncodeToString(chunk)
}

// 无符号64位整数
type Uint64 uint64

func (n Uint64) Encode(v interface{}) []byte {
	chunk := make([]byte, 8)
	binary.BigEndian.PutUint64(chunk, v.(uint64))
	return chunk
}

func (n Uint64) Decode(chunk []byte) interface{} {
	if chunk != nil {
		chunk = common.ResizeBytes(chunk, true, 8)
		return binary.BigEndian.Uint64(chunk)
	}
	return uint64(0)
}

// 无符号Double Word
type Uint32 uint32

func (n Uint32) Encode(v interface{}) []byte {
	chunk := make([]byte, 4)
	binary.BigEndian.PutUint32(chunk, v.(uint32))
	return chunk
}

func (n Uint32) Decode(chunk []byte) interface{} {
	if chunk != nil {
		chunk = common.ResizeBytes(chunk, true, 4)
		return binary.BigEndian.Uint32(chunk)
	}
	return uint32(0)
}

// 无符号24位整数
type Uint24 uint32

func (n Uint24) Encode(v interface{}) []byte {
	chunk := make([]byte, 4)
	binary.BigEndian.PutUint32(chunk, v.(uint32))
	return common.ResizeBytes(chunk, true, 3)
}

func (n Uint24) Decode(chunk []byte) interface{} {
	if chunk != nil {
		chunk = common.ResizeBytes(chunk, true, 4)
		chunk[0] = 0x00 // 去掉最前面一个字节
		return binary.BigEndian.Uint32(chunk)
	}
	return uint32(0)
}

// 无符号Word
type Uint16 uint16

func (n Uint16) Encode(v interface{}) []byte {
	chunk := make([]byte, 2)
	binary.BigEndian.PutUint16(chunk, v.(uint16))
	return chunk
}

func (n Uint16) Decode(chunk []byte) interface{} {
	if chunk != nil {
		chunk = common.ResizeBytes(chunk, true, 2)
		return binary.BigEndian.Uint16(chunk)
	}
	return uint16(0)
}
