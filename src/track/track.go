package track

import (
	"net"

	"github.com/gofrs/uuid/v5"
)

type Track interface {
	Open() error
	Close() error

	// Track under hood:
	// - args
	// - timestamp
	// - response algorhitm
	Event(
		protocol string,
		server_port int,
		client_ip net.IP,
		client_port int,
		chain_id uuid.UUID,
		event_type string,
		payload_size int,
		payload_data []byte,
	) error
}
