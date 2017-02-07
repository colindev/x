SECRETFILE = keys/GCETEST-832b38975e55.json

build:
	go build -o bin/token `ls token/*.go | grep -v _test.go`

cloud-platform:
	cat scopes/cloud-platform | xargs bin/token -secret-file $(SECRETFILE) > token.cache
