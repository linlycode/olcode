package hubpkg

import (
	"io"

	"github.com/linlycode/olcode/pkg/idgen"
)

// Peer denotes the peer entity in a hub
type Peer struct {
	ID     int64
	sender io.Writer
}

// NewPeerID is used to generate unique peerID
func NewPeerID() (int64, error) {
	return idgen.GetIDGen().GenID(idgen.PeerID)
}

// NewPeer create a peer given id
func NewPeer(id int64, sender io.Writer) *Peer {
	return &Peer{
		ID:     id,
		sender: sender,
	}
}

// Send sends message to this peer
func (p *Peer) Send(msg []byte) error {
	_, err := p.sender.Write(msg)
	return err
}
