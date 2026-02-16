package replay

import "net"

type Replay interface {
	Init() error
	Replay(
		client_payload []byte,
		protocol string,
		port int,
		clientIP net.IP,
		clientPort int,
	) []byte
}
