package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP established connection.
type TCPPeer struct {
	net.Conn      // The underlying connection of the peer (TCP conn in this case)
	outbound bool // True if the connection is outbound, false if inbound
	wg       *sync.WaitGroup
}

// NewTCPPeer creates a new TCPPeer.
func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		wg:       &sync.WaitGroup{},
	}
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func (p *TCPPeer) CloseStream() {
	p.wg.Done()
}

// TCPTransportOpts contains the options for the TCP transport.
type TCPTransportOpts struct {
	ListenAddr    string
	Decoder       Decoder
	HandshakeFunk HandshakeFunk
	OnPeer        func(Peer) error
}

// TCPTransport represents the transport layer for TCP connections.
type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC
}

// NewTCPTransport creates a new TCPTransport.
func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC, 1024),
	}
}

// Consume returns a read-only channel for reading incoming messages from peers.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

// ListenAndAccept starts listening for and accepting incoming connections.
func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()
	log.Printf("TCP Transport Listening on port %s\n", t.ListenAddr)
	return nil
}

// Addr implements the transport interface, returning the address the transport is accepting connections.
func (t *TCPTransport) Addr() string {
	return t.ListenAddr
}

// startAcceptLoop continuously accepts incoming connections.
func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			log.Printf("TCP accept error %s\n", err)
			continue
		}

		go t.handleConn(conn, false)
	}
}

// Dial initiates a connection to the specified address.
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)
	return nil
}

// handleConn handles a new connection, performing the handshake and reading messages.
func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error
	defer func() {
		fmt.Printf("dropping peer connection: %s\n", err)
		conn.Close()
	}()
	peer := NewTCPPeer(conn, outbound)

	if err = t.HandshakeFunk(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	// Read Loop
	for {
		rpc := RPC{}
		err = t.Decoder.Decode(conn, &rpc)
		if err != nil {
			fmt.Printf("TCP error: %s\n", err)
			return
		}

		rpc.From = conn.RemoteAddr().String()

		if rpc.Stream {
			peer.wg.Add(1)
			log.Printf("[%s] incoming stream, waiting...\n", conn.RemoteAddr())
			peer.wg.Wait()
			log.Printf("[%s] stream closed, resuming read loop...\n", conn.RemoteAddr())
			continue
		}
		t.rpcch <- rpc
	}
}

// Close closes the transport listener.
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}
