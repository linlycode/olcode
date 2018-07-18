package idgen

import (
	"math/rand"
)

// IDType is the type of ID
type IDType int

const (
	// HubID denotes Hub ID type
	HubID IDType = iota
	// PeerID denotes Peer ID type
	PeerID
)

// IDGen takes responsibility for generating a unique id
type IDGen interface {
	GenID(t IDType) (int64, error)
}

type idGen struct{}

func (g *idGen) GenID(t IDType) (int64, error) {
	return rand.Int63(), nil
}

// single instance of id generator
var singleIDGen = &idGen{}

// GetIDGen query the signle instance IDGen
func GetIDGen() IDGen {
	return singleIDGen
}
