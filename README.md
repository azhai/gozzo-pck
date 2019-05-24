# gozzo 尜舟

## 用途1：Redis请求解析
```go
package main
import (
    "fmt"
    "bytes"
    "strconv"
    "gituhb.com/azhai/gozzo-pck/match"
)

func CreateMatcher(chunk []byte) *match.FieldMatcher {
    offset := bytes.IndexByte(chunk, byte('$'))
    index := bytes.Index(chunk[offset:], []byte("\r\n"))
    length, _ := strconv.Atoi(string(chunk[offset+1:offset+index]))
    m := match.NewFieldMatcher()
	m.AddFixeds([]int{1, offset - 3, 2, 1, index - 1, 2}, nil)
	m.AddField("cmd", match.NewField(size, false))
    return m
}

fun main() {
    chunk := []byte("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n")
    matcher := CreateMatcher(chunk)
    data := matcher.Match(chunk, true)
    if cmd, ok := data["cmd"]; ok {
        fmt.Println("Command is ", string(cmd))
    }
}
```

## 用途2：手机号码归属地
* （可选）生成数据文件 city.txt 和 phone.txt
```bash
cd phone/data/
# 下载原始数据文件 phones.csv （21.3MB） 放入此目录
gawk -f conv.awk phones.csv
cd ../..
```
* 生成可执行文件和数据文件
```bash
cd phone/
chmod +x build.sh && ./build.sh
chmod +x phone && ./phone 1599955 1990123
# 将可执行文件 phone 和数据文件（保留data目录） data/phone.dat 一起打包即可
cd ..
```