package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/azhai/gozzo-pck/find"
	"github.com/azhai/gozzo-utils/common"
	"github.com/azhai/gozzo-utils/filesystem"
)

type ISP = uint8

const (
	MOB_DAT_FILE = "phone.dat"
	MOB_KEY_SIZE = 4
	MOB_POS_SIZE = 2
)
const (
	Unknow ISP = uint8(iota)
	ChnTelecom
	ChnMobile
	ChnUnicom
)

var (
	IspArray = []string{"未知", "中国电信", "中国移动", "中国联通"}
	finder   = NewMobFinder("data/")
)

func Mob2Bin(phone, ispName string) []byte {
	switch ispName {
	default:
		phone += string(Unknow + '0')
	case IspArray[int(ChnTelecom)]:
		phone += string(ChnTelecom + '0')
	case IspArray[int(ChnMobile)]:
		phone += string(ChnMobile + '0')
	case IspArray[int(ChnUnicom)]:
		phone += string(ChnUnicom + '0')
	}
	data, _ := hex.DecodeString(phone)
	return data
}

func Phone2Bin(phone string) []byte {
	if size := len(phone); size < 7 {
		phone = phone + strings.Repeat("0", 7-size)
	}
	phone = phone[:7] + "9" //最大运营商代码
	return common.Hex2Bin(phone)
}

func ReadMobFile(fdir string) (records []string, keypairs []find.KeyPair, err error) {
	// 较小的文件使用ReadLines
	fpath := filepath.Join(fdir, "city.txt")
	lines, err := filesystem.ReadLines(fpath)
	if err != nil {
		return
	}
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		ps := strings.SplitN(line, "\t", 2)
		if len(ps) < 2 {
			continue
		}
		records = append(records, ps[1])
	}
	// 较大的文件使用LineReader
	fpath = filepath.Join(fdir, "phone.txt")
	reader := filesystem.NewLineReader(fpath)
	defer reader.Close()
	if reader.Err() != nil {
		return
	}
	for reader.Reading() {
		line := reader.Text()
		if len(line) == 0 {
			continue
		}
		ps := strings.SplitN(line, "\t", 3)
		if len(ps) < 3 {
			continue
		}
		n, _ := strconv.Atoi(ps[2])
		key := Mob2Bin(ps[0], ps[1])
		keypairs = append(keypairs, find.KeyPair{Key: key, Idx: n})
	}
	return
}

type MobFinder struct {
	*find.Finder
}

func NewMobFinder(fdir string) *MobFinder {
	fpath := filepath.Join(fdir, MOB_DAT_FILE)
	fp, size, err := filesystem.OpenFile(fpath, false, false)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if size <= 0 {
		records, keypairs, _ := ReadMobFile(fdir)
		builder := find.NewBuilder(MOB_KEY_SIZE, MOB_POS_SIZE)
		builder.Build(fp, records, keypairs)
	}
	f := new(MobFinder)
	f.Finder, _ = find.NewFinder(fp, MOB_KEY_SIZE, MOB_POS_SIZE)
	return f
}

func (f *MobFinder) Find(phone string) (area, isp string, err error) {
	target := Phone2Bin(phone)
	key, addr := f.SearchIndex(target)
	if addr == nil {
		err = errors.New("Not found")
		return
	}
	data, err := f.GetRecord(addr)
	if err == nil {
		area = string(data)
		x := uint(key[len(key)-1] & 0x0f)
		isp = IspArray[int(x)]
	}
	return
}

func main() {
	var (
		phone string
		size  = len(os.Args)
	)
	if size <= 1 {
		return
	}
	for i := 1; i < size; i++ {
		phone = os.Args[i]
		area, isp, err := finder.Find(phone)
		if err != nil {
			fmt.Errorf("没有找到数据")
		}
		fmt.Println(phone, isp)
		fmt.Println(area)
	}
}
