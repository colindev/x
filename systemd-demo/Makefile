build:
	go build -o demo-process ./main.go

install:
	cp ./demo-process /usr/local/sbin/
	cp ./demo.service /lib/systemd/system/demo.service

run:
	systemctl stop demo
	systemctl daemon-reload && systemctl start demo && systemctl status demo

