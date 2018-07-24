package hubpkg

import (
	"io"

	"github.com/linlycode/olcode/pkg/idgen"
)

var peerIDGen = &idgen.IDGenerator{}

// Peer denotes the peer entity in a hub
type Peer struct {
	ID     int64
	sender io.Writer
}

// NewPeer create a peer given id
func NewPeer(sender io.Writer) *Peer {
	return &Peer{
		ID:     peerIDGen.GenerateID(),
		sender: sender,
	}
}

// Send sends message to this peer
func (p *Peer) Send(msg []byte) error {
	_, err := p.sender.Write(msg)
	return err
}
