package udp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"time"

	"github.com/hezhijie/tiny_transport/util"
	"github.com/pion/dtls/v2"
	"github.com/pion/transport/v2/udp"
)

type Side int

const (
	Client Side = 0
	Server Side = 1
)

type Conn struct {
	tls          bool
	conn         net.Conn
	remoteAddr   net.Addr
	side         Side
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewClientSideUdpConnWrap(conn net.Conn, tls bool) *Conn {
	return &Conn{
		tls:  tls,
		conn: conn,
		side: Client,
	}
}

func NewServerSideUdpConnWrap(conn net.Conn, tls bool) *Conn {
	return &Conn{
		tls:        tls,
		conn:       conn,
		side:       Server,
		remoteAddr: conn.RemoteAddr(),
	}
}

//func (wrap *Conn) SetReadTimeout(readTimeout time.Duration) {
//	wrap.readTimeout = readTimeout
//}
//
//func (wrap *Conn) SetWriteTimeout(writeTimeout time.Duration) {
//	wrap.writeTimeout = writeTimeout
//}

func NewServerSideUdpConnWrapWithRemote(conn *udp.Conn, remoteAddr net.Addr, tls bool) *Conn {
	return &Conn{
		tls:        tls,
		conn:       conn,
		side:       Server,
		remoteAddr: remoteAddr,
	}
}

func ListenTls(port int, pubKey string, crt string, cacrt string) (net.Listener, error) {
	addrStr := fmt.Sprintf("%s:%d", "", port) // 启动服务不指定本地ip
	addr, err := net.ResolveUDPAddr("udp", addrStr)
	certificate, err := util.LoadKeyAndCertificate(pubKey, crt)
	util.Check(err)

	rootCertificate, err := util.LoadCertificate(cacrt)
	util.Check(err)
	certPool := x509.NewCertPool()
	cert, err := x509.ParseCertificate(rootCertificate.Certificate[0])
	util.Check(err)
	certPool.AddCert(cert)

	// Prepare the configuration of the DTLS connection
	config := &dtls.Config{
		Certificates:         []tls.Certificate{certificate},
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
		ClientAuth:           dtls.RequireAndVerifyClientCert,
		ClientCAs:            certPool,
		ConnectContextMaker: func() (context.Context, func()) {
			return context.Background(), func() {
				// do nothing
			}
		},
	}
	return dtls.Listen("udp", addr, config)
}

func Listen(port int) (net.Listener, error) {
	addrStr := fmt.Sprintf("%s:%d", "", port) // 启动服务不指定本地ip
	addr, _ := net.ResolveUDPAddr("udp", addrStr)
	return udp.Listen("udp", addr)
}

func Dial(ip string, port int) (*Conn, error) {
	addr := fmt.Sprintf("%s:%d", ip, port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	conn := NewClientSideUdpConnWrap(udpConn, false)
	return conn, nil
}

func DialTls(ip string, port int, pubKey string, crt string, cacrt string) (*Conn, error) {
	addr := &net.UDPAddr{IP: net.ParseIP(ip), Port: port}
	certificate, err := util.LoadKeyAndCertificate(pubKey, crt)
	util.Check(err)
	rootCertificate, err := util.LoadCertificate(cacrt)
	util.Check(err)
	certPool := x509.NewCertPool()
	cert, err := x509.ParseCertificate(rootCertificate.Certificate[0])
	util.Check(err)
	certPool.AddCert(cert)
	config := &dtls.Config{
		Certificates:         []tls.Certificate{certificate},
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
		RootCAs:              certPool,
	}
	udpConn, err := dtls.DialWithContext(context.Background(), "udp", addr, config)
	if err != nil {
		return nil, err
	}
	return NewClientSideUdpConnWrap(udpConn, false), nil
}

func (wrap *Conn) Read(data []byte) (int, error) {
	if wrap.tls {
		return wrap.read0(data)
	}

	if wrap.side == Server {
		n, err := wrap.read0(data)
		if wrap.remoteAddr == nil {
			wrap.remoteAddr = wrap.conn.(*udp.Conn).RemoteAddr()
		}
		return n, err
	}

	if wrap.side == Client {
		return wrap.read0(data)
	}
	panic("on read, unknown udp type")

}

func (wrap *Conn) Write(data []byte) (int, error) {
	if wrap.tls {
		return wrap.conn.Write(data)
	}
	if wrap.side == Server {
		if wrap.remoteAddr == nil {
			panic("server side ,but remoteAddr is nil")
		}
		return wrap.conn.Write(data)
	}
	if wrap.side == Client {
		return wrap.conn.Write(data)
	}
	panic("on write, unknown udp type")
}

func (wrap *Conn) RemoteAddr() net.Addr {
	return wrap.remoteAddr
}

func (wrap *Conn) read0(data []byte) (int, error) {
	return wrap.conn.Read(data)
}

func (wrap *Conn) LocalAddr() net.Addr {
	return wrap.conn.LocalAddr()
}

func (wrap *Conn) Close() error {
	return wrap.conn.Close()
}
