package main

import (
	"time"
	"bufio"
	"encoding/binary"
	"io"
	"fmt"
	"os"
	"log"
	"bytes"
)

type Package struct {
	Version        [2]byte // 协议版本       //  固定2个字节
	//Length         int16   // 数据部分长度    // 固定的2个字节 ;可以表示 2^16=32k长度 一般足够用了;不够可以使用四个字节
	Length         int32   // 数据部分长度    // 固定的4个字节 ;可以表示 2^32=4g 长度 足够用了！


	// 以下是变长; payload!  (8+ 2+len(Hostname) + 2+len(TagLength) + len(Msg))  ===  Length字段的实际的值！
	Timestamp      int64   // 时间戳

	HostnameLength int16   // 主机名长度
	Hostname       []byte  // 主机名

	TagLength      int16   // Tag长度
	Tag            []byte  // Tag

	Msg            []byte  // 数据部分长度
}

//在网络中是大端序 ！
func (p *Package) Pack(writer io.Writer) error {
	var err error
	err = binary.Write(writer, binary.BigEndian, &p.Version)
	err = binary.Write(writer, binary.BigEndian, &p.Length)
	err = binary.Write(writer, binary.BigEndian, &p.Timestamp)
	err = binary.Write(writer, binary.BigEndian, &p.HostnameLength)
	err = binary.Write(writer, binary.BigEndian, &p.Hostname)
	err = binary.Write(writer, binary.BigEndian, &p.TagLength)
	err = binary.Write(writer, binary.BigEndian, &p.Tag)
	err = binary.Write(writer, binary.BigEndian, &p.Msg)
	return err
}

func (p *Package) Unpack(reader io.Reader) error {
	var err error
	err = binary.Read(reader, binary.BigEndian, &p.Version)
	err = binary.Read(reader, binary.BigEndian, &p.Length)

	err = binary.Read(reader, binary.BigEndian, &p.Timestamp)
	err = binary.Read(reader, binary.BigEndian, &p.HostnameLength)
	p.Hostname = make([]byte, p.HostnameLength)
	err = binary.Read(reader, binary.BigEndian, &p.Hostname)

	err = binary.Read(reader, binary.BigEndian, &p.TagLength)
	p.Tag = make([]byte, p.TagLength)
	err = binary.Read(reader, binary.BigEndian, &p.Tag)

	p.Msg = make([]byte, p.Length-8-2-int32(p.HostnameLength)-2-int32(p.TagLength))
	err = binary.Read(reader, binary.BigEndian, &p.Msg)
	return err
}

func (p *Package) String() string {
	return fmt.Sprintf("version:%s length:%d timestamp:%d hostname:%s tag:%s msg:%s",
		p.Version,
		p.Length,
		p.Timestamp,
		p.Hostname,
		p.Tag,
		p.Msg,
	)
}

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	pack := &Package{
		Version:        [2]byte{'V', '0'},
		Timestamp:      time.Now().Unix(),
		HostnameLength: int16(len(hostname)),
		Hostname:       []byte(hostname),
		TagLength:      4,
		Tag:            []byte("demo"),
		Msg:            []byte(("现在时间是:" + time.Now().Format("2006-01-02 15:04:05"))),
	}
	pack.Length = 8 + 2 + int32(pack.HostnameLength) + 2 + int32(pack.TagLength) + int32(len(pack.Msg))

	buf := new(bytes.Buffer)
	// 写入四次，模拟TCP粘包效果
	pack.Pack(buf)
	pack.Pack(buf)
	pack.Pack(buf)
	pack.Pack(buf)
	// scanner
	scanner := bufio.NewScanner(buf)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if !atEOF && data[0] == 'V' {
			if len(data) > 6 {
				length := int32(0)
				binary.Read(bytes.NewReader(data[2:6]), binary.BigEndian, &length)
				if int(length)+ 6 <= len(data) {
					return int(length) + 6, data[:int(length)+6], nil
				}
			}
		}
		return
	})
	for scanner.Scan() {
		scannedPack := new(Package)
		scannedPack.Unpack(bytes.NewReader(scanner.Bytes()))
		log.Println(scannedPack)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("无效数据包")
	}
}
