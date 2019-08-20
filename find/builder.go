package find

import (
	"bytes"
	"io"
	"time"

	"github.com/azhai/gozzo-pck/serialize"
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
	IdxBegin  uint32
	IdxEnd  uint32
	KeyCount  uint32
	SizeProps uint16
	Version string
	*serialize.Object
}

func NewDatHeader(keySize, positSize int) *DatHeader {
	p := &DatHeader{
		KeySize: keySize,
		PositSize: positSize,
		Object: serialize.NewObject(),
	}
	p.AddUintField("idxBegin", 4) // 第一个索引开始位置
	p.AddUintField("idxEnd", 4)   // 最后一个索引结束位置
	p.AddUintField("keyCount", 4)
	// 0-7 ItemSize: 0（变长）~ 256
	// 8-12 KeySize: 1 ~ 31
	// 13-15 PositSize: 2 ~ 4
	p.AddUintField("sizeProps", 2)
	p.AddHexStrField("version", 4) // 4字节
	return p
}

func (h *DatHeader) GetIndexRange() (int, int) {
	return int(h.IdxBegin), int(h.IdxEnd)
}

func (h *DatHeader) GetHeaderSize() int {
	return FIX_BYTES*3 + 2 + VER_BYTES
	_, size := h.Matcher.GetLeastSize()
	return size
}

func (h *DatHeader) GetSizeProps(itemSize int) uint16 {
	return uint16(h.PositSize<<13 + h.KeySize<<8 + itemSize)
}

type DatIndex struct {
	Key []byte
	Pos uint64
	*serialize.Object
}

func NewDatIndex(keySize, positSize int) *DatIndex {
	p := &DatIndex{Object: serialize.NewObject()}
	p.AddBytesField("key", keySize, false)
	p.AddUintField("pos", positSize)
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
	b.Header.SizeProps = b.Header.GetSizeProps(0)
	base := b.Header.GetHeaderSize()
	if b.Header.IdxBegin, err = b.BuildRecord(rs, base); err != nil {
		return err
	}
	if b.Header.KeyCount, err = b.BuildIndex(ks); err != nil {
		return err
	}
	b.Header.IdxEnd = b.Header.IdxBegin + uint32(b.Index.Len())
	b.Header.Version = time.Now().Format("060102") + "00"
	headBytes := serialize.Serialize(b.Header)
	if _, err = w.Write(headBytes); err != nil {
		return err
	}
	if _, err := b.Record.WriteTo(w); err != nil {
		return err
	}
	_, err = b.Index.WriteTo(w)
	return
}

func (b *Builder) BuildRecord(records []string, base int) (idxBegin uint32, err error) {
	var pos Position
	for _, rec := range records {
		pos = Position(base + b.Record.Len())
		_, err = b.Record.WriteString(rec)
		b.Record.WriteByte(0x00)
		b.PosList = append(b.PosList, pos)
	}
	idxBegin = uint32(base + b.Record.Len()) // IdxBegin
	return
}

func (b *Builder) BuildIndex(keypairs []KeyPair) (keyCount uint32, err error) {
	var addr Position
	positCount := len(b.PosList)
	for _, pair := range keypairs {
		if pair.Idx < 0 || pair.Idx >= positCount {
			addr = Position(0)
		} else {
			addr = b.PosList[pair.Idx]
		}
		b.IdxObject.Key = pair.Key
		b.IdxObject.Pos = uint64(addr)
		chunk := serialize.Serialize(b.IdxObject)
		_, err = b.Index.Write(chunk)
	}
	if size := len(keypairs); size > 0 {
		keyCount = uint32(size)
	}
	return
}
