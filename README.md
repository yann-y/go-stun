go-stun
=======

go-stun is a STUN (RFC 3489, 5389) client implementation in golang
(a.k.a. UDP hole punching).

[RFC 3489](https://tools.ietf.org/html/rfc3489):
STUN - Simple Traversal of User Datagram Protocol (UDP)
Through Network Address Translators (NATs)

[RFC 5389](https://tools.ietf.org/html/rfc5389):
Session Traversal Utilities for NAT (STUN)

### Use the Command Line Tool

Simply run these commands (if you have installed golang and set `$GOPATH`)
```
go get github.com/yann-y/go-stun
go-stun
```
or clone this repo and run these commands
```
make build
./bin/go-stun
```
You will get the output like
```
NAT Type: Port restricted NAT
External IP Family: 1
External IP: 110.191.219.11
External Port: 5671
```
You can use `-s` flag to use another STUN server, and use `-v` to work on
verbose mode.
```bash
> ./bin/go-stun --help
Usage of ./bin/go-stun:
  -debug
        double verbose mode (includes -info)
  -info
        verbose mode
  -loop int
        loop interval (company: s)
  -p int
        port to listen on for client (default 13333)
  -s string
        STUN server address (default "stun.qq.com:3478")
  -ttl
        loop interval (company: s)
```

### Use the Library

The library `github.com/yann-y/go-stun/stun` is extremely easy to use -- just
one line of code.

```go
import "github.com/yann-y/go-stun/stun"

func main() {
	nat, host, err := stun.NewClient().Discover()
}
```
