package service

import (
	"errors"
	"net/http"
	"strings"

	"github.com/TadahTech/foodlogiq-demo/pkg/model"
)

// So this is the only thing I questioned, but for the sake of a demo I think it's fine.
// Since we are not expiring tokens, and they are statically assigned, there's not a real point in storing in a DB
// This wouldn't be the same in a prod env but for a demo, it's fine
var (
	acme = &model.User{
		UserID: 12345,
		Name:   "Acme",
		Token:  "74edf612f393b4eb01fbc2c29dd96671",
	}

	ajax = &model.User{
		UserID: 12345,
		Name:   "98765 ",
		Token:  "d88b4b1e77c70ba780b56032db1c259b",
	}
	usersByToken = map[string]*model.User{
		acme.Token: acme,
		ajax.Token: ajax,
	}
)

func userFromBearer(req *http.Request) (*model.User, error) {
	reqToken := req.Header.Get("Authorization")

	if len(reqToken) == 0 {
		return nil, errors.New("no Authorization token")
	}

	splitToken := strings.Split(reqToken, "Bearer ")

	if len(splitToken) != 2 {
		return nil, errors.New("bearer token malformed")
	}

	token := splitToken[1]
	if len(token) == 0 {
		return nil, errors.New("token empty")
	}

	u := usersByToken[token]
	if u == nil {
		return nil, errors.New("no user found for that token")
	}

	return u, nil
}

func tokenMw(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := userFromBearer(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
