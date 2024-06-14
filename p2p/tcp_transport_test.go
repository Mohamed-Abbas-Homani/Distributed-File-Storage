package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOpts {
		ListenAddr: ":3000",
		Decoder: DefaultDecoder{},
		HandshakeFunk: NOPHandshakeFunc,
	}
	tr := NewTCPTransport(opts)

	assert.Equal(t, tr.ListenAddr, opts.ListenAddr)

    assert.Nil(t, tr.ListenAndAccept())

}
