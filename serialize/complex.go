package serialize

import (
	"fmt"
	"time"

	"github.com/azhai/gozzo-utils/common"
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

func NewTwoDim(size, quad int) *TwoDim {
	return &TwoDim{
		Quadrant: quad,
		Unsigned: &Unsigned{Size: size},
	}
}

func NewTwoDimXY(size int, x, y int64) *TwoDim {
	td := NewTwoDim(size, QuadrantFirst)
	if y < 0 {
		y = 0 - y
		td.Quadrant = QuadrantForth
	}
	if x < 0 {
		x = 0 - x
		// y的符号已经去掉，这里不能再用 y<0 来判断
		if td.Quadrant == QuadrantForth {
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
	var dx, dy uint64
	switch v := v.(type) {
	case TwoDim:
		dx, dy = v.Xdim, v.Ydim
	case *TwoDim:
		dx, dy = v.Xdim, v.Ydim
	}
	fmt.Printf("%d %d\n", dx, dy)
	xbs := td.Unsigned.Encode(dx)
	ybs := td.Unsigned.Encode(dy)
	fmt.Printf("%x %x\n", xbs, ybs)
	return append(xbs, ybs...)
}

func (td *TwoDim) Decode(chunk []byte) interface{} {
	td.Xdim = td.Unsigned.DecodeUint64(chunk[:td.Size])
	td.Ydim = td.Unsigned.DecodeUint64(chunk[td.Size:])
	return td
}

// 时间戳，精确到秒
type TimeStamp struct {
	*Integer
}

func NewTimeStamp() *TimeStamp {
	return &TimeStamp{
		Integer: &Integer{
			Unsigned: &Unsigned{Size: 8},
		},
	}
}

func (ts TimeStamp) Encode(v interface{}) []byte {
	if t, ok := v.(time.Time); ok {
		return ts.Integer.Encode(t.Unix())
	}
	return nil
}

func (ts TimeStamp) Decode(chunk []byte) interface{} {
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
	t, err := common.ParseDate(LayoutDate, v.(string))
	if err == nil {
		return t
	}
	return nil
}