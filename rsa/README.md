# RSA Exercise

### Generator

```golang
./rsa-generate -key /tmp/key.pem -pub /tmp/pub.pem -len 2048
```

### Sign

```golang
echo 'Hello world' | ./rsa-sign -key /tmp/key.pem -sha 256 > /tmp/signed
```

### Verify

```golang
echo 'Hello world' | ./rsa-verify -key /tmp/key.pem -sign /tmp/signed -sha 256 
```
