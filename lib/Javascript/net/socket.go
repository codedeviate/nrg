package jsnet

import "net"

func SocketOpen(address net.Addr) (net.Conn, error) {
	conn, err := net.Dial(address.Network(), address.String())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func SocketClose(connection net.Conn) error {
	return connection.Close()
}

func SocketWrite(connection net.Conn, data string) (int, error) {
	return connection.Write([]byte(data))
}

func SocketRead(connection net.Conn) (string, error) {
	buf := make([]byte, 1024)
	n, err := connection.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[0:n]), nil
}

func SocketListen(address net.Addr) (net.Listener, error) {
	listener, err := net.Listen(address.Network(), address.String())
	if err != nil {
		return nil, err
	}
	return listener, nil
}
