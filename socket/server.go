package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

type collector struct {
	m map[net.Conn]bool
	*sync.RWMutex
}

func (c *collector) add(conn net.Conn) {
	c.Lock()
	defer c.Unlock()

	c.m[conn] = true
}

func (c *collector) del(conn net.Conn) {
	c.Lock()
	defer c.Unlock()

	delete(c.m, conn)
}

func (c *collector) handle(conn net.Conn) {
	c.add(conn)
	log.Println("conn ", len(c.m))
	defer func() {
		c.del(conn)
		log.Println("conn ", len(c.m))
		conn.Close()
	}()

	buf := bufio.NewReader(conn)
	for {
		b, _, err := buf.ReadLine()
		if err != nil {
			if err != io.EOF {
				os.Stderr.WriteString(err.Error())
			}
			return
		}

		os.Stdout.Write(b)
		os.Stdout.WriteString("\n")
	}
}

func (c *collector) each(fn func(net.Conn)) {
	c.RLock()
	defer c.RUnlock()

	for conn := range c.m {
		fn(conn)
	}
}

func listen(network, addr string, args []string) {
	var (
		hub = &collector{
			RWMutex: &sync.RWMutex{},
			m:       map[net.Conn]bool{},
		}
		tcpAddr     *net.TCPAddr
		tcpListener *net.TCPListener
		err         error
	)

	tcpAddr, err = net.ResolveTCPAddr(network, addr)
	checkError(err)

	tcpListener, err = net.ListenTCP(network, tcpAddr)
	checkError(err)

	adaptor := serverAdapt(args)
	go func() {
		adaptor(hub)
		tcpListener.Close()
	}()

	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go hub.handle(conn)
	}
}
