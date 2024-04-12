package middleware

import (
	"errors"
	"goapi/api"
	"goapi/internal/tools"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var unAuthorizedError = errors.New("Invalid username or token")

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var username string = r.URL.Query().Get("username")
		var token = r.Header.Get("Authorization")
		var err error

		if username == "" || token == "" {
			log.Error(unAuthorizedError)
			api.RequestErrorHandler(w, unAuthorizedError)
			return
		}

		var database *tools.DatabaseInterface
		database, err = tools.NewDatabase()

		if err != nil {
			api.InternalErrorHandler(w)
			return
		}

		var loginDetails *tools.LoginDetails = (*database).GetUserLoginDetails(username)

		if loginDetails == nil || (token != (*loginDetails).AuthToken) {
			log.Error(unAuthorizedError)
			api.RequestErrorHandler(w, unAuthorizedError)

			return
		}

		next.ServeHTTP(w, r)
	})
}
