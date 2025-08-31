package main

import (
	"fmt"
	"net/http"

	"github.com/7-Dany/dev/auth"
	"github.com/7-Dany/dev/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) authMiddleware(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error:%v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByAuthKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get user:%v", err))
			return
		}

		handler(w, r, user)
	}
}
