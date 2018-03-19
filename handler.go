package olcode

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type respCode int

const (
	success respCode = iota
	notLoggedIn
	roomNotExist
)

type response struct {
	Code respCode    `json:"code"`
	Data interface{} `json:"data"`
}

type handler struct {
	userStore   *userStore
	roomManager *roomManager
	wsUpgrader  *websocket.Upgrader
	wsMtx       sync.Mutex
	wsConns     map[int64]map[roomID]*websocket.Conn
}

func newHandler() *handler {
	registerSessionTypes()

	return &handler{
		userStore:   newUserStore(),
		roomManager: newRoomManager(),
		wsUpgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		wsConns: make(map[int64]map[roomID]*websocket.Conn),
	}
}

func checkMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		log.Printf("invalid method %s", r.Method)
		http.Error(w, "HTTP method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func parseRequest(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		log.Printf("invalid Content-Type %q", contentType)
		http.Error(w, "invalid Content-Type", http.StatusBadRequest)
		return false
	}

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(v); err != nil {
		log.Printf("failed to decode request body, %s", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return false
	}
	return true
}

func replyServerError(w http.ResponseWriter, reason string) {
	log.Printf(reason)
	http.Error(w, "", http.StatusInternalServerError)
}

func reply(w http.ResponseWriter, code respCode, data interface{}) {
	resp := &response{Code: code, Data: data}
	bytes, err := json.Marshal(resp)
	if err != nil {
		replyServerError(w, fmt.Sprintf("fail to encode response, %s", err))
		return
	}

	if code != success {
		log.Printf("reply with code %v", code)
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write(bytes); err != nil {
		replyServerError(w, fmt.Sprintf("fail to write response, %s", err))
	}
}

func getRequestUser(r *http.Request) (*User, error) {
	session, err := getSession(r)
	if err != nil {
		return nil, fmt.Errorf("fail to get session, %s", err)
	}

	key := "user"
	v, ok := session.Values[key]
	if !ok {
		return nil, nil
	}

	user, ok := v.(*User)
	if !ok {
		return nil, fmt.Errorf("session[%q] is not a *User, value: %v", key, v)
	}
	return user, nil
}

type loginRequest struct {
	Name string `json:"name"`
}

type loginResponse struct {
	UserID int64 `json:"user_id"`
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	log.Printf("login")

	if !checkMethod(w, r, http.MethodPost) {
		return
	}

	var req loginRequest
	if ok := parseRequest(w, r, &req); !ok {
		return
	}

	session, err := getSession(r)
	if err != nil {
		replyServerError(w, fmt.Sprintf("fail to get session, %s", err))
		return
	}

	user := h.userStore.newUser(req.Name)
	session.Values["user"] = user

	if err := session.Save(r, w); err != nil {
		replyServerError(w, fmt.Sprintf("failed to save session, %s", err))
		return
	}

	reply(w, success, loginResponse{UserID: user.ID})
}

type createResponse struct {
	RoomID roomID `json:"room_id"`
}

func (h *handler) create(w http.ResponseWriter, r *http.Request) {
	log.Printf("create")

	if !checkMethod(w, r, http.MethodPost) {
		return
	}

	user, err := getRequestUser(r)
	if err != nil {
		replyServerError(w, fmt.Sprintf("fail to get user from session, %s", err))
		return
	}
	if user == nil {
		reply(w, notLoggedIn, nil)
		return
	}

	id := h.roomManager.create(user)
	resp := &createResponse{RoomID: id}

	reply(w, success, resp)
}

type attendRequest struct {
	RoomID roomID `json:"room_id"`
}

func (h *handler) attend(w http.ResponseWriter, r *http.Request) {
	log.Printf("attend")

	if !checkMethod(w, r, http.MethodGet) {
		return
	}

	vals, ok := r.URL.Query()["room_id"]
	if !ok || len(vals) < 1 {
		http.Error(w, "missing room_id parameter", http.StatusBadRequest)
		return
	}
	id := roomID(vals[0])

	user, err := getRequestUser(r)
	if err != nil {
		replyServerError(w, fmt.Sprintf("fail to get user from session, %s", err))
		return
	}
	if user == nil {
		reply(w, notLoggedIn, nil)
		return
	}

	if err := h.roomManager.attend(id, user); err != nil {
		reply(w, roomNotExist, nil)
		return
	}

	h.wsMtx.Lock()
	defer h.wsMtx.Unlock()

	userConns, ok := h.wsConns[user.ID]
	if !ok {
		h.wsConns[user.ID] = make(map[roomID]*websocket.Conn)
	}

	if _, ok := userConns[id]; !ok {
		conn, err := h.wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			replyServerError(w, fmt.Sprintf("failed to upgrade websocket, %s", err))
		}
		userConns[id] = conn
	}

	// conn.WriteJSON()

	// reply(w, success, nil)
}

type leaveRequest attendRequest

func (h *handler) leave(w http.ResponseWriter, r *http.Request) {
	log.Printf("leave")

	if !checkMethod(w, r, http.MethodPost) {
		return
	}

	var req leaveRequest
	if ok := parseRequest(w, r, &req); !ok {
		return
	}

	user, err := getRequestUser(r)
	if err != nil {
		replyServerError(w, fmt.Sprintf("fail to get user from session, %s", err))
		return
	}
	if user == nil {
		reply(w, notLoggedIn, nil)
		return
	}

	if err := h.roomManager.leave(req.RoomID, user); err != nil {
		reply(w, roomNotExist, nil)
		return
	}

	reply(w, success, nil)
}
