SECRETFILE = keys/GCETEST-832b38975e55.json

build:
	go build -o bin/jwt `ls jwt/*.go | grep -v _test.go`
	go build -o bin/oauth2google `ls oauth2google/*go | grep -v _test.go`

jwt-cloud-platform:
	cat scopes/cloud-platform | xargs bin/jwt -secret-file $(SECRETFILE) > jwt.cache

oauth2-cloud-platform:
	cat scopes/cloud-platform | xargs bin/oauth2google -secret-file $(SECRETFILE) oauth2google.cache
