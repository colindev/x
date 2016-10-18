package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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

	if err := ioutil.WriteFile(pidfile, []byte(strconv.Itoa(os.Getpid())), os.ModePerm); err != nil {
		log.Println(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
		log.Println(r.URL)
	})

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Println(err)
	}

}
