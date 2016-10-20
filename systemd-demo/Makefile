build:
	go build -o demo-process ./main.go

install:
	cp ./demo-process /usr/local/sbin/
	cp ./demo.service /lib/systemd/system/demo.servic
	systemctl enable demo && systemctl daemon-reload && systemctl status demo

