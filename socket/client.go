package main

import (
	"bufio"
	"io"
	"net"
	"os"
)

func dial(network, addr string, args []string) {
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

	var out = clientAdapt(args)

	go func() {
		buf := bufio.NewReader(tcpConn)
		for {
			b, _, err := buf.ReadLine()
			checkError(err)

			out(b)
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
