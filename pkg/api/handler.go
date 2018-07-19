package api

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/linlycode/olcode/pkg/api/wsconn"
	"github.com/linlycode/olcode/pkg/hubpkg"
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
	log.Printf("new connection from: %s\n", r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	ah := wsconn.NewAsyncHandler(conn, h.hm)
	ah.AsyncRun()
}
