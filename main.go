package main

import (
	// "bytes"
	"io"
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
    go s2.Start()
    time.Sleep(2 * time.Second)
	// data := bytes.NewReader([]byte("my file"))
	
  //   s2.Store("mydata", data)
    r, err := s2.Get("mydata")
		if err != nil {
			log.Fatal(err)
		}

		b, err := io.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(string(b))
    select {}
}   
