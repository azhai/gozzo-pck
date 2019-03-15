package tool

import (
	"bytes"
)

/* 掩码Enum
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

type WalkFunc func(item interface{}) error
type MapFunc func(item interface{}) (interface{}, error)
type ReduceFunc func(a, b interface{}) (interface{}, error)

type IArray interface {
	ToList() []interface{}
}

// 循环修改数组
func ArrayWalk(arr IArray, f WalkFunc) error {
	for _, item := range arr.ToList() {
		if err := f(item); err != nil {
			return err
		}
	}
	return nil
}

func ArrayMap(arr IArray, f MapFunc) ([]interface{}, error) {
	var (
		res []interface{}
		err error
	)
	for _, item := range arr.ToList() {
		if r, err := f(item); err != nil {
			res = append(res, r)
		}
	}
	return res, err
}

func ArrayReduce(arr IArray, f ReduceFunc, res interface{}) (interface{}, error) {
	var err error
	for _, item := range arr.ToList() {
		res, err = f(res, item)
	}
	return res, err
}
