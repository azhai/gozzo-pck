package common

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

/*
// 掩码Enum
type Weekday = uint8
const (
	Sunday Weekday = uint8(iota)
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

type Property = uint32
const (
	Bit0 Property = 1 << iota
	Bit1
	_
	Bit3
)
*/

//16进制的字符表示和二进制字节表示互转
var Bin2Hex = hex.EncodeToString

var Hex2Bin = func(data string) []byte {
	if block, err := hex.DecodeString(data); err == nil {
		return block
	}
	return nil
}

//字符求余转为数字
func ToNum(c byte) uint8 {
	return uint8((c - '0' + 100) % 10)
}

//将数值或者数字转为字符串
func ToString(v interface{}) string {
	return fmt.Sprintf("%#v", v)
}

// 补充空字节
func ExtendBytes(data []byte, isLeft bool, size int) []byte {
	if size <= 0 {
		return data
	}
	padding := bytes.Repeat([]byte{0x00}, size)
	if isLeft {
		return append(padding, data...)
	} else {
		return append(data, padding...)
	}
}

// 调整长度
func ResizeBytes(data []byte, isLeft bool, n int) []byte {
	size := len(data) - n // 多余长度
	if size == 0 {
		return data
	} else if size < 0 {
		return ExtendBytes(data, isLeft, 0-size)
	}
	if isLeft {
		return data[size:]
	} else {
		return data[:size]
	}
}
