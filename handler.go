package olcode

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
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
	userStore *userStore
	hubMgr    *HubMgr
	homePath  string
}

func newHandler(homePath string) *handler {
	registerSessionTypes()

	return &handler{
		userStore: newUserStore(),
		hubMgr:    NewHubMgr(),
		homePath:  homePath,
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
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	log.Printf("login")

	session, err := getSession(r)
	if err != nil {
		replyServerError(w, fmt.Sprintf("fail to get session, %s", err))
		return
	}

	if r.Method == http.MethodGet {
		u, ok := session.Values["user"]
		if !ok || u == nil {
			reply(w, notLoggedIn, nil)
			return
		}
		if user, ok := u.(*User); ok {
			if _, ok := h.userStore.users[user.ID]; !ok {
				reply(w, notLoggedIn, nil)

				delete(session.Values, "user")
				if err := session.Save(r, w); err != nil {
					replyServerError(w, fmt.Sprintf("failed to save session, %s", err))
				}
				return
			}
			reply(w, success, loginResponse{UserID: user.ID, Name: user.Name})
		} else {
			replyServerError(w, fmt.Sprintf("user in session is not *User, session[user]: %v", u))
		}
		return
	}

	if !checkMethod(w, r, http.MethodPost) {
		return
	}

	var req loginRequest
	if ok := parseRequest(w, r, &req); !ok {
		return
	}

	user := h.userStore.newUser(req.Name)
	session.Values["user"] = user

	if err := session.Save(r, w); err != nil {
		replyServerError(w, fmt.Sprintf("failed to save session, %s", err))
		return
	}

	reply(w, success, loginResponse{UserID: user.ID, Name: user.Name})
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

	id, err := h.hubMgr.registerHub(user)
	if err != nil {
		replyServerError(w, fmt.Sprintf("fail to register room, err=%v", err))
		return
	}
	reply(w, success, &createResponse{RoomID: id})
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
	if !ok || len(vals) == 0 {
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

	hub, err := h.hubMgr.getHub(id)
	if err != nil {
		if err == errRoomNotExist {
			reply(w, roomNotExist, nil)
		} else {
			replyServerError(w, fmt.Sprintf("fail to query hub, id=%v, err=%v", id, err))
		}
		return
	}

	buildClientRoomConn(user, hub, w, r)
}

func (h *handler) serverHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("query home index")

	if !checkMethod(w, r, http.MethodGet) {
		return
	}
	http.ServeFile(w, r, filepath.Join(h.homePath, "index.html"))
}
