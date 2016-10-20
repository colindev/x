package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {

	var (
		pidfile string
		name    string
		addr    string
	)

	flag.StringVar(&pidfile, "pid", "/tmp/demo-process.pid", "pidfile")
	flag.StringVar(&name, "name", "undefined", "app name")
	flag.StringVar(&addr, "addr", ":8000", "http serve on")
	flag.Parse()

	// 輸出 pid
	if err := ioutil.WriteFile(pidfile, []byte(strconv.Itoa(os.Getpid())), os.ModePerm); err != nil {
		log.Println(err)
	}

	go runHTTP(addr)

	// 處理 signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGQUIT)
	s := <-c
	os.Stdout.WriteString(fmt.Sprintf("%v(%#v)\n", s, s))
}

func runHTTP(addr string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
		log.Println(r.URL)
	})

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println(err)
	}
}
