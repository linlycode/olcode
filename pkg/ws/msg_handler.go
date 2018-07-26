package ws

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/linlycode/olcode/pkg/common"
	"github.com/linlycode/olcode/pkg/hubpkg"
)

type msgHandler struct {
	hub    *hubpkg.Hub
	peer   *hubpkg.Peer
	hm     hubpkg.HubMgr
	sender io.WriteCloser
}

// ackHello acknowledges the hello command from the client
// command format: ACKHELLO {HUB_ID}
func (h *msgHandler) ackHello(success bool) error {
	successInt := 0
	if success {
		successInt = 1
	}
	_, err := h.sender.Write([]byte(fmt.Sprintf("ACKHELLO %d %d %d", successInt, h.hub.ID, h.peer.ID)))
	return err
}

func (h *msgHandler) notifyJoined() {
	common.Assert(h.hub != nil)
	if err := h.hub.Broadcast(h.peer.ID, []byte(fmt.Sprintf("PEER_JOINED"))); err != nil {
		log.WithError(err).Error("fail to broadcast, sender will be closed")
		h.sender.Close()
	}
}

func (h *msgHandler) Handle(msg []byte) (handleErr error) {
	msgS := string(msg)
	switch {
	case strings.HasPrefix(msgS, "HELLO"):
		success := false
		defer func() {
			handleErr = h.ackHello(success)
		}()

		if h.hub != nil {
			return fmt.Errorf("only one HELLO command is allowed")
		}

		tokens := strings.Split(msgS, " ")
		common.Assertf(len(tokens) <= 2 && len(tokens) >= 1, "invalid HELLO command: %s", msgS)

		var hub *hubpkg.Hub
		var peer *hubpkg.Peer
		var err error
		if len(tokens) == 2 {
			// got command: HELLO {HUB_ID}
			hubID, err := strconv.ParseInt(tokens[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid hub id: %s", tokens[1])
			}
			if hub = h.hm.GetHub(hubID); hub == nil {
				return fmt.Errorf("hub#%d not exist", hubID)
			}
		} else if hub, err = h.hm.GenHub(); err != nil {
			// got command: HELLO
			return fmt.Errorf("fail to generate hub, err=%v", err)
		}

		peer = hubpkg.NewPeer(h.sender)
		if err := hub.AddPeer(peer); err != nil {
			return fmt.Errorf("fail to add peer to hub, err=%d", err)
		}

		success = true
		h.peer = peer
		h.hub = hub
		go h.notifyJoined()
		return

	default:
		// all the other message just broadcast
		common.Assert(h.hub != nil)
		return h.hub.Broadcast(h.peer.ID, msg)
	}
}
