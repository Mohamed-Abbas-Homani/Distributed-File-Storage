package main

import (
	// "bytes"
	"bytes"
	"fmt"
	// "io"
	"log"
	"time"

	"github.com/Mohamed-Abbas-Homani/dfs/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunk: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)
	opts := FileServerOpts{

		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(opts)
	tcpTransport.OnPeer = s.OnPeer
	return s
}

func main() {

	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")
	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(2 * time.Second)
	go func() {
		err := s2.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(2 * time.Second)

	for i := 0; i < 10; i++ {
		data := bytes.NewReader([]byte("my file"))
		err := s2.Store(fmt.Sprintf("mydata_%d", i), data)
		if err != nil {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
	// r, err := s2.Get("my data")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// b, err := io.ReadAll(r)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println(string(b))
	select {}
}
