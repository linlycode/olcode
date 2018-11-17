package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type contextKey struct {
	key string
}

var (
	contextJSONDataKey = contextKey{"json_data_key"}
	// ErrCtxJSONDataNotFound defines the error for json data key not found error
	ErrCtxJSONDataNotFound = errors.New("json data key not found")
)

func jsonMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			if !strings.HasPrefix(req.RequestURI, "/api") {
				next.ServeHTTP(resp, req)
				return
			}

			body := req.Body
			if body != nil {
				defer body.Close()
				jData, err := ioutil.ReadAll(body)
				if err != nil {
					log.WithError(err).Error("fail read body")
					resp.WriteHeader(400)
					return
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, contextJSONDataKey, jData)
				req = req.WithContext(ctx)
			}

			next.ServeHTTP(resp, req)
		})
}

func getJSONData(req *http.Request) ([]byte, error) {
	data, ok := req.Context().Value(contextJSONDataKey).([]byte)
	if !ok {
		return nil, ErrCtxJSONDataNotFound
	}
	return data, nil
}

func h200(resp http.ResponseWriter, data interface{}) {
	rawData, err := json.Marshal(data)
	if err != nil {
		log.WithError(err).Error("fail to marshal data")
	}
	resp.Write(rawData)
}

func h400(resp http.ResponseWriter, err error) {
	log.WithError(err).Error("internal error")
	resp.WriteHeader(400)
	resp.Write([]byte(fmt.Sprintf("%s", err)))
}

func h500(resp http.ResponseWriter, err error) {
	log.WithError(err).Error("internal error")
	resp.WriteHeader(500)
}
