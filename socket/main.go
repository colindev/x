package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		listenMode bool
		network    string
		addr       string
	)

	flag.BoolVar(&listenMode, "l", false, "use listen mode")
	flag.StringVar(&network, "net", "tcp", "network")
	flag.StringVar(&addr, "addr", "127.0.0.1:8000", "address")
	flag.Parse()

	if listenMode {
		listen(network, addr)
	} else {
		dial(network, addr)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, syscall.SIGHUP)
	for {
		log.Println("[signal] ", <-s)
	}
}

func listen(network, addr string) {
	var (
		tcpAddr     *net.TCPAddr
		tcpListener *net.TCPListener
		err         error
	)

	tcpAddr, err = net.ResolveTCPAddr(network, addr)
	checkError(err)

	tcpListener, err = net.ListenTCP(network, tcpAddr)
	checkError(err)

	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handler(conn)
	}
}

func handler(conn net.Conn) {

	defer conn.Close()

	buf := bufio.NewReader(conn)
	for {
		b, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				log.Println(conn, err)
			}
			return
		}

		fmt.Println(conn, string(b))
	}
}

func dial(network, addr string) {
	var (
		tcpAddr *net.TCPAddr
		tcpConn io.ReadWriteCloser
		err     error
	)
	tcpAddr, err = net.ResolveTCPAddr(network, addr)
	checkError(err)

	tcpConn, err = net.DialTCP(network, nil, tcpAddr)
	checkError(err)
	defer tcpConn.Close()

	go func() {
		buf := bufio.NewReader(tcpConn)
		for {
			b, _, err := buf.ReadLine()
			checkError(err)

			fmt.Println(string(b))
		}
	}()

	buf := bufio.NewReader(os.Stdin)
	for {
		b, _, err := buf.ReadLine()
		checkError(err)

		tcpConn.Write(b)
		tcpConn.Write([]byte("\n"))
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
