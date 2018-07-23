package hubpkg

import (
	"sync"

	"github.com/linlycode/olcode/pkg/common"

	"github.com/linlycode/olcode/pkg/idgen"
)

// HubMgr manages all the hubs
type HubMgr interface {
	GenHub() (*Hub, error)
	GetHub(id int64) *Hub
}

var hubIDGen = &idgen.IDGenerator{}

type hubMgr struct {
	hm  map[int64]*Hub
	mtx *sync.RWMutex
}

func (h *hubMgr) GenHub() (*Hub, error) {
	hID := hubIDGen.GenerateID()

	h.mtx.Lock()
	defer h.mtx.Unlock()
	_, exist := h.hm[hID]
	common.Assert(!exist)

	hub := newHub(hID)
	h.hm[hID] = hub
	return hub, nil
}

func (h *hubMgr) GetHub(id int64) *Hub {
	h.mtx.RLock()
	h.mtx.RUnlock()
	return h.hm[id]
}

// NewHubMgr makes new hub manager
func NewHubMgr() HubMgr {
	return &hubMgr{
		hm:  make(map[int64]*Hub),
		mtx: &sync.RWMutex{},
	}
}
