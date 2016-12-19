package nettools

import (
	"fmt"
	"io"
	"log"
	"net"
)

// Tunnel creates a duplex TCP tunnel between a local and a remote host:port
// pairs. It uses the optional function auth to verify that the tunnel should
// be established based on the incoming connection details. If auth returns
// false, the incoming connection is immediately closed.
func Tunnel(localAddr, remoteAddr string, auth func(conn net.Conn) bool) error {
	local, err := net.Listen("tcp", localAddr)
	if err != nil {
		return fmt.Errorf("creatTunnel: %v", err)
	}
	log.Printf("Creating tunnel %v => %v", localAddr, remoteAddr)
	go func() {
		for {
			conn, err := local.Accept()
			if conn == nil {
				log.Fatalf("accept failed: %v", err)
			}
			if auth == nil || auth(conn) {
				go forward(conn, remoteAddr)
			} else {
				conn.Close()
			}
		}
	}()
	return nil
}

func forward(local net.Conn, remoteAddr string) {
	remote, err := net.Dial("tcp", remoteAddr)
	if remote == nil {
		log.Println("remote dial failed: ", err)
		return
	}
	log.Printf("established %v => %v (%v)", local.RemoteAddr(), remoteAddr, remote.RemoteAddr())
	go func() {
		// Read from 'local' and write to 'remote'. The return values
		// aren't useful here so they are just discarded.
		io.Copy(remote, local)
		// Apparently io.Copy() will exit only when there are problems
		// reading from the source, and not when writes fail. This
		// means that the other duplex channel will not be aware that
		// the remote end disconnected and may keep running
		// indefinitely. Closing the remote endpoint here will do what
		// I need - abort the other io.Copy().
		remote.Close()
	}()
	go func() {
		io.Copy(local, remote)
		local.Close()
	}()
}
