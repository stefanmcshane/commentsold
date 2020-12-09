package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stefanmcshane/commentsold/auth"
)

func authenticateUser(signingMachine auth.JWTToken, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type userAuth struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		ctx := r.Context()

		var ua userAuth
		err := json.NewDecoder(r.Body).Decode(&ua)
		if err != nil {
			logErrorAndReturn(ctx, w, "Must supply username and password in body", http.StatusBadRequest)
			return
		}

		if ua.Username == "commentsold" && ua.Password == "supersecurepassword" {
			jwtAuth := auth.Auth{Username: ua.Username, Password: ua.Password}
			token, err := signingMachine.GenerateToken(ctx, jwtAuth)
			if err != nil {
				logErrorAndReturn(ctx, w, "Unable to generate token", http.StatusServiceUnavailable)
				return
			}
			w.Write([]byte(fmt.Sprintf(`{"token":"%s"}`, token.Token)))
			next.ServeHTTP(w, r)
			return
		}
		logErrorAndReturn(ctx, w, "Incorrect username or password", http.StatusUnauthorized)
	})
}
func checkTokenValidity(signingMachine auth.JWTToken, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		ctx := r.Context()

		st := auth.SignedToken{Token: token}
		_, err := signingMachine.ValidateToken(ctx, st)
		if err != nil {
			e := fmt.Sprintf("Error validating token - %s", err.Error())
			lg(ctx).Errorf(e)
			http.Error(w, e, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
