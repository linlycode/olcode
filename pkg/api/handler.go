package api

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/linlycode/olcode/pkg/hubpkg"
	"github.com/linlycode/olcode/pkg/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024,
	WriteBufferSize: 1024 * 1024,
}

type handler interface {
	// TODO: use a decorator function to process the http Response Writer & Rquest
	serveWS(http.ResponseWriter, *http.Request)
}

type h struct {
	hm hubpkg.HubMgr
}

func newHandler() handler {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	return &h{
		hm: hubpkg.NewHubMgr(),
	}
}

func (h *h) serveWS(w http.ResponseWriter, r *http.Request) {
	log.Infof("new connection from: %s", r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithError(err).Error("fail to upgrage http request")
		return
	}
	ah := ws.NewAsyncHandler(conn, h.hm)
	ah.AsyncRun()
}
