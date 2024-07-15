package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var (
	timeout int64
	size    int
	count   int
	typ uint8 = 8
	code uint8 = 0
)

type ICMP struct{
	Type uint8
	Code uint8
	CheckSum uint16
	ID uint16 
	SequenceNum uint16 
}

func main() {
	GetCommandArgs()
	desIP := os.Args[len(os.Args)-1]
	conn, err := net.DialTimeout("ip:icmp", desIP, time.Duration(timeout)*time.Millisecond)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	remoteAddr := conn.RemoteAddr()
	fmt.Println(remoteAddr)
	for i:=0; i<count; i++{
		icmp := &ICMP{
			Type: typ,
			Code: code,
			CheckSum: 0,
			ID: uint16(i),
			SequenceNum: uint16(i),
		}
		var buffer bytes.Buffer
		binary.Write(&buffer, binary.BigEndian, icmp)
		data := make([]byte, size)
		buffer.Write(data)
		data = buffer.Bytes()
		checkSum := checkSum(data)
		data[2] = byte(checkSum >> 8)
		data[3] = byte(checkSum)

		startTime := time.Now()
		conn.SetDeadline(time.Now().Add(time.Duration(timeout)*time.Millisecond))
		n, err := conn.Write(data)
		if err != nil {
			log.Println(err)
			break
		}
		buf := make([]byte, 1024)
		n, err = conn.Read(buf)
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Printf("来自 %d.%d.%d.%d 的回复：字节=%d 时间=%dms TTL=%d\n", buf[12], buf[13], buf[14], buf[15], n-28, time.Since(startTime).Milliseconds(), buf[8])
		time.Sleep(time.Second)
	}
}

func checkSum(data []byte) uint16{
	length := len(data)
	index := 0
	var sum uint32
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		length -= 2
		index += 2
	}
	if length == 1 {
		sum += uint32(data[index])
	}
	hi := sum >> 16
	for hi != 0 {
		sum = hi + uint32(uint16(sum))
		hi = sum >> 16
	}
	return uint16(^sum)
}

func GetCommandArgs() {
	flag.Int64Var(&timeout, "w", 1000, "请求超时时间")
	flag.IntVar(&size, "l", 32, "发送字节数")
	flag.IntVar(&count, "n", 4, "请求次数")
	flag.Parse()
}