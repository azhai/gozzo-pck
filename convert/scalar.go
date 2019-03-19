package convert

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/azhai/gozzo-utils/common"
)

type IConvert interface {
	Encode() []byte
	Decode(chunk []byte) error
	GetData() interface{}
	SetData(data interface{})
}

// 字节数组
type Bytes struct {
	Data []byte
}

func (s *Bytes) Encode() []byte {
	return s.Data
}

func (s *Bytes) Decode(chunk []byte) error {
	s.Data = chunk
	return nil
}

func (s *Bytes) GetData() interface{} {
	return s.Data
}

func (s *Bytes) SetData(data interface{}) {
	s.Data = data.([]byte)
}

// 字符串
type String struct {
	Data string
}

func (s *String) Encode() []byte {
	return []byte(s.Data)
}

func (s *String) Decode(chunk []byte) error {
	s.Data = string(chunk)
	return nil
}

func (s *String) GetData() interface{} {
	return s.Data
}

func (s *String) SetData(data interface{}) {
	s.Data = data.(string)
}

// BCD码
type HexStr struct {
	Data string
}

func (s *HexStr) Encode() []byte {
	if result, err := hex.DecodeString(s.Data); err == nil {
		return result
	}
	return nil
}

func (s *HexStr) Decode(chunk []byte) error {
	s.Data = hex.EncodeToString(chunk)
	return nil
}

func (s *HexStr) GetData() interface{} {
	return s.Data
}

func (s *HexStr) SetData(data interface{}) {
	s.Data = data.(string)
}

// 无符号64位整数
type Uint64 struct {
	Data uint64
}

func (n *Uint64) Encode() []byte {
	result := make([]byte, 8)
	binary.BigEndian.PutUint64(result, n.Data)
	return result
}

func (n *Uint64) Decode(chunk []byte) error {
	if chunk != nil {
		chunk = common.ResizeBytes(chunk, true, 8)
		n.Data = binary.BigEndian.Uint64(chunk)
	}
	return nil
}

func (n *Uint64) GetData() interface{} {
	return n.Data
}

func (n *Uint64) SetData(data interface{}) {
	n.Data = data.(uint64)
}

// 无符号Double Word
type Uint32 struct {
	Data uint32
}

func (n *Uint32) Encode() []byte {
	result := make([]byte, 4)
	binary.BigEndian.PutUint32(result, n.Data)
	return result
}

func (n *Uint32) Decode(chunk []byte) error {
	if chunk != nil {
		chunk = common.ResizeBytes(chunk, true, 4)
		n.Data = binary.BigEndian.Uint32(chunk)
	}
	return nil
}

func (n *Uint32) GetData() interface{} {
	return n.Data
}

func (n *Uint32) SetData(data interface{}) {
	n.Data = data.(uint32)
}

// 无符号24位整数
type Uint24 struct {
	Uint32
}

func (n *Uint24) Encode() []byte {
	return common.ResizeBytes(n.Uint32.Encode(), true, 3)
}

// 无符号Word
type Uint16 struct {
	Data uint16
}

func (n *Uint16) Encode() []byte {
	result := make([]byte, 2)
	binary.BigEndian.PutUint16(result, n.Data)
	return result
}

func (n *Uint16) Decode(chunk []byte) error {
	if chunk != nil {
		chunk = common.ResizeBytes(chunk, true, 2)
		n.Data = binary.BigEndian.Uint16(chunk)
	}
	return nil
}

func (n *Uint16) GetData() interface{} {
	return n.Data
}

func (n *Uint16) SetData(data interface{}) {
	n.Data = data.(uint16)
}
