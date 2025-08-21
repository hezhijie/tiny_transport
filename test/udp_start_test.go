package test

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hezhijie/tiny_transport/example"
	"github.com/hezhijie/tiny_transport/udp"
)

type TestHandler struct {
}

func (t *TestHandler) OnRead(conn *udp.Conn, data []byte) {
	log.Printf("onRead, remote:%s,data:%s", conn.RemoteAddr(), string(data))
}

func (t *TestHandler) OnClose(conn *udp.Conn) {
	//TODO implement me
	panic("implement me")
}

func TestStart(t *testing.T) {
	example.StartUdp(&TestHandler{})
	select {}
}

func TestDial(t *testing.T) {
	addr := fmt.Sprintf("%s:%d", "127.0.0.1", 40002)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Printf("ResolveUDPAddr exp,%v", err)
		return
	}
	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	for i := 0; i < 10; i++ {
		data := []byte("Hello " + strconv.Itoa(i))
		_, err := udpConn.Write(data)
		if err != nil {
			fmt.Printf("Write exp,%v", err)
			break
		}
		time.Sleep(time.Second * 1)
	}
	err = udpConn.Close()
	if err != nil {
		fmt.Printf("Close exp,%v", err)
	}
	//select {}
}

func TestDial2(t *testing.T) {

	udpConn, err := udp.Dial("127.0.0.1", 40002)

	//udpConn, err := net.DialUDP("udp", nil, udpAddr)
	for i := 0; i < 10; i++ {
		data := []byte("World " + strconv.Itoa(i))
		_, err := udpConn.Write(data)
		if err != nil {
			fmt.Printf("Write exp,%v", err)
			break
		}
		time.Sleep(time.Second * 1)
	}
	err = udpConn.Close()
	if err != nil {
		fmt.Printf("Close exp,%v", err)
	}
	//select {}
}

func TestStartTls(t *testing.T) {
	example.StartUdpTls(&TestHandler{})
	select {}
}

func TestDialTls(t *testing.T) {

	pubKey := "/Users/hezhijie/GolandProjects/tiny_transport/example/client.pem"
	crt := "/Users/hezhijie/GolandProjects/tiny_transport/example/client.crt"
	cacrt := "/Users/hezhijie/GolandProjects/tiny_transport/example/server.crt"

	udpConn, err := udp.DialTls("127.0.0.1", 40002, pubKey, crt, cacrt)

	//udpConn, err := net.DialUDP("udp", nil, udpAddr)
	for i := 0; i < 10; i++ {
		data := []byte("World " + strconv.Itoa(i))
		_, err := udpConn.Write(data)
		if err != nil {
			fmt.Printf("Write exp,%v", err)
			break
		}
		time.Sleep(time.Second * 1)
	}
	//err = udpConn.Close()
	if err != nil {
		fmt.Printf("Close exp,%v", err)
	}
	//select {}
}

func TestTimeout(t *testing.T) {
	num := getsth()
	log.Printf("get sth num:%d", num)
}

func getsth() int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	var ret int
	go func() {
		defer func() {
			log.Printf("cancel...")
			cancel()
		}()
		time.Sleep(time.Second * 2)
		log.Printf("after sleep..")
		ret = 5
	}()

	for {
		select {
		case <-ctx.Done():
			if ctx.Err() != nil && strings.Contains(ctx.Err().Error(), "context deadline exceeded") {
				log.Printf("timeout...,%v", ctx.Err())
				ret = -1
			}
			return ret
		default:
		}
	}
}
