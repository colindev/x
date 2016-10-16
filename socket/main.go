package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	var (
		network string
		in      string
		out     string
		args    []string
	)

	flag.StringVar(&in, "in", "", "stream in")
	flag.StringVar(&out, "out", "", "stream out")
	flag.StringVar(&network, "net", "tcp", "network")
	flag.Parse()

	args = flag.Args()

	if out != "" {
		listen(network, out, args)
	}

	if in != "" {
		dial(network, in, args)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGTERM, syscall.SIGHUP)
	for {
		log.Println("[signal] ", <-s)
	}
}
