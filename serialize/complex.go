package serialize

import (
	"time"
)

const (
	// 平面坐标系的四个象限
	QuadrantFirst  = 1 // 第一象限，全正
	QuadrantSecond = 2 // 第二象限，x负y正
	QuadrantThird  = 3 // 第二象限，全负
	QuadrantForth  = 4 // 第二象限，x正y负
	// 默认日期格式
	LayoutDate = "20060102"
)

// 二维数据(双整型)，可用于复数、平面中的点、经纬度坐标（先转为Decimal）
type TwoDim struct {
	Xdim, Ydim uint64
	Quadrant int
	*Unsigned
}

func NewTwoDim(quad int) *TwoDim {
	return &TwoDim{
		Quadrant: quad,
		Unsigned: &Unsigned{Size: 8},
	}
}

func NewTwoDimXY(x, y int64) *TwoDim {
	td := NewTwoDim(QuadrantFirst)
	if y < 0 {
		y = 0 - y
		td.Quadrant = QuadrantForth
	}
	if x < 0 {
		x = 0 - x
		if y < 0 {
			td.Quadrant = QuadrantThird
		} else {
			td.Quadrant = QuadrantSecond
		}
	}
	td.Xdim = uint64(x)
	td.Ydim = uint64(y)
	return td
}

func (td TwoDim) Encode(v interface{}) []byte {
	x := td.Unsigned.Encode(td.Xdim)
	y := td.Unsigned.Encode(td.Ydim)
	return append(x, y...)
}

func (td *TwoDim) Decode(chunk []byte) interface{} {
	v := NewTwoDim(td.Quadrant)
	v.Xdim = td.Unsigned.Decode(chunk[:td.Size]).(uint64)
	v.Ydim = td.Unsigned.Decode(chunk[td.Size:]).(uint64)
	return v
}

// 时间戳，精确到秒
type TimeStamp struct {
	*Integer
}

func (ts TimeStamp) Encode(v interface{}) []byte {
	if t, ok := v.(time.Time); ok {
		return ts.Integer.Encode(t.Unix())
	}
	return nil
}

func (ts *TimeStamp) Decode(chunk []byte) interface{} {
	v := ts.Integer.DecodeInt64(chunk)
	return time.Unix(v, 0)
}

// 日期，格式20060102
type Date struct {
	HexStr
}

func (d Date) Encode(v interface{}) []byte {
	if t, ok := v.(time.Time); ok {
		return d.HexStr.Encode(t.Format(LayoutDate))
	}
	return nil
}

func (d Date) Decode(chunk []byte) interface{} {
	v := d.HexStr.Decode(chunk)
	if t, err := time.Parse(LayoutDate, v.(string)); err == nil {
		return t
	}
	return nil
}