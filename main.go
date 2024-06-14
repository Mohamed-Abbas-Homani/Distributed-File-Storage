package main

import (
	"fmt"
	"log"

	"github.com/Mohamed-Abbas-Homani/dfs/p2p"
)

func OnPeer(p2p.Peer) error {
	fmt.Printf("doing some logic witht the peer outside of TCPTransport\n")
	return nil
 }
func main() {
	opts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		Decoder:       p2p.DefaultDecoder{},
		HandshakeFunk: p2p.NOPHandshakeFunc,
		OnPeer:        OnPeer,
	}
	tr := p2p.NewTCPTransport(opts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("Message: %+v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}
}
