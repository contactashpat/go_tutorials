package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	addr := flag.String("addr", ":9000", "address to listen on, e.g., ':9000' or '127.0.0.1:7000'")
	flag.Parse()

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to listen on %s: %v\n", *addr, err)
		os.Exit(1)
	}
	defer listener.Close()

	log.Printf("Echo server listening on %s", *addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	log.Printf("client connected: %s", conn.RemoteAddr())
	if _, err := io.Copy(conn, conn); err != nil {
		log.Printf("io error for %s: %v", conn.RemoteAddr(), err)
		return
	}
	log.Printf("client disconnected: %s", conn.RemoteAddr())
}
