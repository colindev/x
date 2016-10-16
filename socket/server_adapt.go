package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
)

func serverAdapt(args []string) func(hub *collector) {

	if len(args) < 1 {
		return func(hub *collector) {
			buf := bufio.NewReader(os.Stdin)
			for {
				b, _, err := buf.ReadLine()
				if err != nil {
					if err != io.EOF {
						log.Println(err)
					}
					msg := []byte("server down: " + err.Error())
					hub.each(func(conn net.Conn) {
						conn.Write(msg)
					})

					return
				}

				hub.each(func(conn net.Conn) {
					_, err := conn.Write(b)
					conn.Write([]byte("\n"))
					if err != nil {
						log.Printf("write %v failed: %v\n", conn, err)
					}
				})
			}
		}
	}

	return func(hub *collector) {
		cmd := exec.Command(args[0], args[1:]...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}

		buf := bufio.NewReader(stdout)
		for {
			b, err := buf.ReadBytes('\n')
			if err != nil {
				if err != io.EOF {
					log.Println("[error]", err)
				}
				return
			}

			hub.each(func(conn net.Conn) {
				_, err := conn.Write(b)
				if err != nil {
					log.Printf("write %v failed: %v\n", conn, err)
				}
			})
		}

	}
}
