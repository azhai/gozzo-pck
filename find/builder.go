package find

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/azhai/gozzo-pck/convert"
)

// 格式说明
//  ---------------------------
// |  Header 头部
//  ---------------------------
// |  Record 记录区
//  ---------------------------
// |  Subset 子索引，可选
//  ---------------------------
// |  Index  顶层索引
//  ---------------------------

const (
	ITEM_SIZE_MAX = 256 // ItemSize的最大值
	FIX_BYTES     = 4   // IdxBegin、IdxEnd、KeyCount的字节数
	VER_BYTES     = 4   // 版本号的字节数
)

type Position = uint32

type KeyPair struct {
	Key []byte
	Idx int
}

type DatHeader struct {
	KeySize   int
	PositSize int
	convert.Object
}

func NewDatHeader(keySize, positSize int) *DatHeader {
	p := &DatHeader{KeySize: keySize, PositSize: positSize}
	p.Init()
	p.AddUint32Field("idxBegin") // 第一个索引开始位置
	p.AddUint32Field("idxEnd")   // 最后一个索引结束位置
	p.AddUint32Field("keyCount")
	// 0-7 ItemSize: 0（变长）~ 256
	// 8-12 KeySize: 1 ~ 31
	// 13-15 PositSize: 2 ~ 4
	p.AddUint16Field("sizeProps")
	p.AddHexStrField("version", 8) // 4字节
	return p
}

func (h *DatHeader) GetIndexRange() (int, int) {
	return int(h.GetUint32("idxBegin")), int(h.GetUint32("idxEnd"))
}

func (h *DatHeader) GetHeaderSize() int {
	return FIX_BYTES*3 + 2 + VER_BYTES
}

func (h *DatHeader) GetSizeProps(itemSize int) uint16 {
	return uint16(h.PositSize<<13 + h.KeySize<<8 + itemSize)
}

type DatIndex struct {
	convert.Object
}

func NewDatIndex(keySize, positSize int) *DatIndex {
	p := &DatIndex{}
	p.Init()
	p.AddBytesField("key", keySize)
	if positSize == 2 {
		p.AddUint16Field("pos")
	} else if positSize == 3 {
		p.AddUint24Field("pos")
	} else {
		p.AddUint32Field("pos")
	}
	return p
}

type Builder struct {
	Header    *DatHeader
	Record    bytes.Buffer
	Index     bytes.Buffer
	IdxObject *DatIndex
	PosList   []Position
}

func NewBuilder(keySize, positSize int) *Builder {
	return &Builder{
		Header:    NewDatHeader(keySize, positSize),
		IdxObject: NewDatIndex(keySize, positSize),
	}
}

func (b *Builder) Build(w io.Writer, rs []string, ks []KeyPair) (err error) {
	var idxBegin Position
	headData := make(map[string]interface{})
	headData["sizeProps"] = b.Header.GetSizeProps(0)
	if idxBegin, err = b.BuildRecord(rs); err != nil {
		return err
	}
	if headData["keyCount"], err = b.BuildIndex(ks); err != nil {
		return err
	}
	idxEnd := idxBegin + Position(b.Index.Len())
	fmt.Println(idxBegin, idxEnd)
	headData["idxBegin"] = uint32(idxBegin)
	headData["idxEnd"] = uint32(idxEnd)
	version := time.Now().Format("060102")
	headData["version"] = version + "00"
	b.Header.SetTable(headData)
	headBytes := b.Header.Encode()
	if _, err = w.Write(headBytes); err != nil {
		return err
	}
	if _, err := b.Record.WriteTo(w); err != nil {
		return err
	}
	_, err = b.Index.WriteTo(w)
	return
}

func (b *Builder) BuildRecord(records []string) (idxBegin Position, err error) {
	var pos Position
	base := b.Header.GetHeaderSize()
	for _, rec := range records {
		pos = Position(base + b.Record.Len())
		_, err = b.Record.WriteString(rec)
		b.Record.WriteByte(0x00)
		b.PosList = append(b.PosList, pos)
	}
	idxBegin = Position(base + b.Record.Len()) // IdxBegin
	return
}

func (b *Builder) BuildIndex(keypairs []KeyPair) (keyCount uint32, err error) {
	var addr Position
	positCount := len(b.PosList)
	positSize := b.Header.PositSize
	idxData := make(map[string]interface{})
	for _, pair := range keypairs {
		idxData["key"] = pair.Key
		if pair.Idx < 0 || pair.Idx >= positCount {
			addr = Position(0)
		} else {
			addr = b.PosList[pair.Idx]
		}
		if positSize == 2 {
			idxData["pos"] = uint16(addr)
		} else {
			idxData["pos"] = uint32(addr)
		}
		b.IdxObject.SetTable(idxData)
		_, err = b.Index.Write(b.IdxObject.Encode())
	}
	if size := len(keypairs); size > 0 {
		keyCount = uint32(size)
	}
	return
}
