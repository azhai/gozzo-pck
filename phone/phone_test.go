package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindPhone(t *testing.T) {
	phones := []string{
		"15999558910123123213213",
		"1300",
		"1703576",
		"199997922323"}
	for _, phone := range phones {
		area, isp, err := finder.Find(phone)
		assert.NoError(t, err)
		t.Log(phone, isp)
		t.Log(area)
	}
}

func TestFindError(t *testing.T) {
	_, _, err := finder.Find("afsd32323")
	assert.Error(t, err, "错误的结果")
	t.Log(err)
}

func BenchmarkFindPhone(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		var i = 0
		for p.Next() {
			i++
			phone := fmt.Sprintf("%s%d%s", "1897", i&10000, "45")
			_, _, err := finder.Find(phone)
			assert.NoError(b, err)
		}
	})
}
