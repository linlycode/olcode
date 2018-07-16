package apiservice

import (
	"net/http"

	"github.com/linly/olcode/pkg/wshub"
)

type handler interface {
	// TODO: use a decorator function to process the http Response Writer & Rquest
	createHub(http.ResponseWriter, *http.Request)
	joinHub(http.ResponseWriter, *http.Request)
}

type h struct {
	hm wshub.HubMgr
}

func newHandler() handler {
	return &h{}
}

func (*h) createHub(http.ResponseWriter, *http.Request) {}
func (*h) joinHub(http.ResponseWriter, *http.Request)   {}
