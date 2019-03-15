package main

import (
	"fmt"
	"testing"
)

func FindPhoneTest(t *testing.T, phone string) {
	area, isp, err := finder.Find(phone)
	if err != nil {
		t.Fatal("没有找到数据")
	}
	t.Log(phone, isp)
	t.Log(area)
}

func TestFindPhone(t *testing.T) {
	FindPhoneTest(t, "15999558910123123213213")
	FindPhoneTest(t, "1300")
	FindPhoneTest(t, "1703576")
	FindPhoneTest(t, "199997922323")
}

func TestFindError(t *testing.T) {
	_, _, err := finder.Find("afsd32323")
	if err == nil {
		t.Fatal("错误的结果")
	}
	t.Log(err)
}

func BenchmarkFindPhone(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		var i = 0
		for p.Next() {
			i++
			phone := fmt.Sprintf("%s%d%s", "1897", i&10000, "45")
			_, _, err := finder.Find(phone)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
