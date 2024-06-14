package p2p

// HandshakeFunk... ?
type HandshakeFunk func(any) error

func NOPHandshakeFunc(any) error {return nil}