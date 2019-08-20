package find

import (
	"bytes"
	"io"
	"sort"

	"github.com/azhai/gozzo-pck/serialize"
)

func GetAddrUint32(addr []byte) uint32 {
	obj := serialize.NewUnsigned(4)
	return obj.Decode(addr).(uint32)
}

type ReaderIndex interface {
	SetUnitSize(size int)
	Count() int
	ReadIndex(data []byte, idx int) (int, error)
}

// 折半查找
func BinSearch(r ReaderIndex, target []byte, checkTop bool) int {
	key := make([]byte, len(target))
	count := r.Count()
	find := func(i int) bool {
		r.ReadIndex(key, i)
		return bytes.Compare(target, key) < 0
	}
	if find(0) {
		return -1 // 超出下限
	} else if checkTop && !find(-1) {
		return -2 // 超出上限
	}
	return sort.Search(count, find) - 1
}

// 索引等长的目录
type Catalog struct {
	data     []byte
	length   int
	count    int
	unitSize int
}

func NewCatalog(data []byte) *Catalog {
	c := &Catalog{length: len(data)}
	if c.length > 0 {
		c.data = make([]byte, c.length)
		copy(c.data[:], data[:])
	}
	return c
}

func (c *Catalog) SetSource(reader io.ReaderAt, offset, size int) error {
	c.data = make([]byte, size)
	n, err := reader.ReadAt(c.data, int64(offset))
	if err == nil {
		c.length = n
	}
	return err
}

func (c *Catalog) SetUnitSize(size int) {
	c.unitSize = size
	if c.unitSize > 0 {
		c.count = c.length / c.unitSize
	}
}

func (c *Catalog) Count() int {
	return c.count
}

func (c *Catalog) ReadAt(data []byte, offset int64) (int, error) {
	start := int(offset)
	stop := start + len(data)
	return copy(data[:], c.data[start:stop]), nil
}

func (c *Catalog) ReadIndex(data []byte, idx int) (int, error) {
	if idx < 0 {
		idx += c.count
	}
	offset := int64(c.unitSize * idx)
	return c.ReadAt(data, offset)
}

type Finder struct {
	reader  io.ReaderAt
	catalog *Catalog
	Header  *DatHeader
}

func NewFinder(reader io.ReaderAt, keySize, positSize int) (f *Finder, err error) {
	f = &Finder{
		reader:  reader,
		catalog: NewCatalog(nil),
		Header:  NewDatHeader(keySize, positSize),
	}
	headData := make([]byte, f.Header.GetHeaderSize())
	if _, err = f.reader.ReadAt(headData, 0); err != nil {
		return
	}
	if err = serialize.Unserialize(headData, f.Header); err != nil {
		return
	}
	idxBegin, idxEnd := f.Header.GetIndexRange()
	f.catalog.SetSource(f.reader, idxBegin, idxEnd-idxBegin)
	return
}

func (f *Finder) SearchIndex(target []byte) ([]byte, []byte) {
	unitSize := f.Header.KeySize + f.Header.PositSize
	f.catalog.SetUnitSize(unitSize)
	i := BinSearch(f.catalog, target, true)
	if i < 0 {
		return nil, nil
	}
	index := make([]byte, unitSize)
	f.catalog.ReadIndex(index, i)
	sep := f.Header.KeySize
	return index[:sep], index[sep:]
}

func (f *Finder) GetRecord(addr []byte) ([]byte, error) {
	data := make([]byte, ITEM_SIZE_MAX)
	offset := int64(GetAddrUint32(addr))
	_, err := f.reader.ReadAt(data, offset)
	if err == nil {
		n := bytes.IndexByte(data, 0x00)
		if n >= 0 {
			data = data[:n]
		}
	}
	return data, err
}
