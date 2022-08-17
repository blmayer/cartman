package cartman

import (
	"crypto/tls"
	"fmt"
	"net"
	"testing"
)

func TestMain(t *testing.T) {
	cert, err := tls.LoadX509KeyPair("fullchain.pem", "privkey.pem")
	if err != nil {
		t.Error(err)
	}
	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		ClientAuth:   tls.RequireAnyClientCert,
		InsecureSkipVerify: true,
	}

	// cartman part
	Debug = true
	cartman, err := NewFileStore("users")
	if err != nil {
		t.Error(err)
	}

	// server listen
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	println("server ready on port 8080")

	for {
		tcp, err := listener.Accept()
		if err != nil {
			println("accept error:", err.Error())
			continue
		}
		println("got a connection")

		// select certificate
		conn := tls.Server(tcp, tlsCfg)
		err = conn.Handshake()
		if err != nil {
			println("handshake error:", err.Error())
			conn.Close()
			continue
		}
		connInfo := conn.ConnectionState()
		fmt.Printf("got request from %+v\n", connInfo)
		println("fingerprint:", fingerprint(connInfo.PeerCertificates[0]))

		userName, err := cartman.GetClientFromCert(connInfo.PeerCertificates[0])
		if err != nil {
			println("user not found")
			t.Error(err)
		}
		println("connected user:", userName)
		conn.Write([]byte("OK\r\n"))
	}
}
