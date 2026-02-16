package tracks

import (
	"log"
	"net"
	"perdoon/src/state"
	"perdoon/src/track"

	"github.com/gofrs/uuid/v5"
)

type Print struct{}

func NewPrint(state *state.State) track.Track {
	return &Print{}
}

func (p *Print) Open() error {
	return nil
}
func (p *Print) Close() error {
	return nil
}

func (p *Print) Event(
	protocol string,
	server_port int,
	client_ip net.IP,
	client_port int,
	chain_id uuid.UUID,
	event_type string,
	payload_size int,
	payload_data []byte,
) error {
	log.Printf("protocol: %s, server_port: %d, client_ip: %s, client_port: %d, chain_id: %s, event_name: %s, payload_size: %d", protocol, server_port, client_ip, client_port, chain_id, event_type, payload_size)

	return nil
}
