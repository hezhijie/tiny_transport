package udp

type ConnHandler interface {
	OnRead(conn *Conn, data []byte)
	OnClose(conn *Conn)
}
