package olcode

import (
	"encoding/gob"
	"net/http"

	"github.com/gorilla/sessions"
)

var sessionStore = sessions.NewCookieStore([]byte("@The Secret #928347"))

func getSession(r *http.Request) (*sessions.Session, error) {
	session, err := sessionStore.Get(r, "olcode")
	if err != nil {
		return nil, err
	}
	return session, nil
}

func registerSessionTypes() {
	gob.Register(&User{})
}
