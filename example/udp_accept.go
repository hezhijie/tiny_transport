package example

import (
	"log"

	"github.com/hezhijie/tiny_transport/udp"
)

func StartUdp(handler udp.ConnHandler) {

	listener, _ := udp.Listen(40002)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("accept exp:%v", err)
			continue
		}
		log.Print("accept a new connection,remote:", conn.RemoteAddr())

		wrapConn := udp.NewServerSideUdpConnWrap(conn, false)

		go func(wrapConn *udp.Conn) {
			for {
				data := make([]byte, 1500)
				l, _ := wrapConn.Read(data)
				handler.OnRead(wrapConn, data[:l])
			}
		}(wrapConn)
	}
}

func StartUdpTls(handler udp.ConnHandler) {

	pubKey := "/Users/hezhijie/GolandProjects/tiny_transport/example/server.pem"
	crt := "/Users/hezhijie/GolandProjects/tiny_transport/example/server.crt"
	cacrt := "/Users/hezhijie/GolandProjects/tiny_transport/example/server.crt"

	listener, _ := udp.ListenTls(40002, pubKey, crt, cacrt)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("accept exp:%v", err)
			continue
		}
		log.Print("accept a new connection,remote:", conn.RemoteAddr())

		wrapConn := udp.NewServerSideUdpConnWrap(conn, false)

		go func(wrapConn *udp.Conn) {
			for {
				data := make([]byte, 1500)
				l, _ := wrapConn.Read(data)
				if l == 0 {
					log.Print("read 0 byte data..")
					continue
				}
				handler.OnRead(wrapConn, data[:l])
			}
		}(wrapConn)
	}
}
