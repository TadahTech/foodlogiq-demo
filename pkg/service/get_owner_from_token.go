package service

import (
	"errors"
	"github.com/TadahTech/foodlogiq-demo/pkg/model"
	"net/http"
	"strings"
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
	user = []*model.User{
		acme, ajax,
	}
	usersByToken = map[string]*model.User{
		acme.Token: acme,
		ajax.Token: ajax,
	}
)

func userFromBearer(req *http.Request) (*model.User, error) {
	reqToken := req.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")

	if len(splitToken) != 2 {
		return nil, errors.New("bearer token malformed")
	}

	token := splitToken[1]
	if len(token) == 0 {
		return nil, errors.New("token empty")
	}

	user := usersByToken[token]
	if user == nil {
		return nil, errors.New("no user found for that token")
	}

	return user, nil
}
