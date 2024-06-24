package p2p

const (
	IncomingStream  = 0x2
	IncomingMessage = 0x1
)

// RPC holds any arbitrary data that is being sent over the
// each transport between two nodes in the network.
type RPC struct {
	From    string
	Payload []byte
	Stream  bool
}
