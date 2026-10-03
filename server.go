package main

import (
    "context"
    "crypto/rand"
    "crypto/ed25519"
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "io"
    "math/big"
    "os"
    "sync"
    "time"

    "github.com/quic-go/quic-go"
)

func main() {
    //Disable auto warning about increasing UDP buffer size
    os.Setenv("QUIC_GO_DISABLE_RECEIVE_BUFFER_WARNING", "true")

    //Configure QUIC keep-alives and higher idle timeout
    quicConfig := &quic.Config {
        MaxIdleTimeout:  120 * time.Second,
        KeepAlivePeriod: 10 * time.Second,
    }

    listener, err := quic.ListenAddr("0.0.0.0:" + os.Args[1], generateTLSConfig(), quicConfig)
    if err != nil {
        fmt.Printf("Listen error: %v\n", err)
        return
    }
    defer listener.Close()

    fmt.Printf("*Listening on UDP %s...\n", os.Args[1])

    ctx := context.Background()
    conn, err := listener.Accept(ctx)
    if err != nil {
        fmt.Printf("Accept error: %v\n", err)
        return
    }
    defer conn.CloseWithError(0, "")

    fmt.Printf("Connection from %s\n", conn.RemoteAddr().String())

    stream, err := conn.AcceptStream(ctx)
    if err != nil {
        fmt.Printf("Stream error: %v\n", err)
        return
    }
    defer stream.Close()

    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        io.Copy(os.Stdout, stream)
    }()

    go func() {
        defer wg.Done()
        io.Copy(stream, os.Stdin)
    }()

    wg.Wait()
    fmt.Println("\nConnection closed.")
}

//Make tls cert. Proto string on the server and client nned to be the same
func generateTLSConfig() *tls.Config {
    pub, priv, _ := ed25519.GenerateKey(rand.Reader)
    template := x509.Certificate{SerialNumber: big.NewInt(1)}
    certDER, _ := x509.CreateCertificate(rand.Reader, &template, &template, pub, priv)

    tlsCert := tls.Certificate{
        Certificate: [][]byte{certDER},
        PrivateKey:  priv,
    }

    return &tls.Config{
        Certificates: []tls.Certificate{tlsCert},
        NextProtos:   []string{"abracadabra"},
    }
}
