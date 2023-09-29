package main

import (
	"fmt"
	"net/http"

	"github.com/witczakxd/first-web-api/internal/auth"
	"github.com/witczakxd/first-web-api/internal/database"
)

type authHandler func(http.ResponseWriter,*http.Request,database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey,err := auth.GetAPIKey(r.Header)
		if err != nil {
			responseWithError(w,403,fmt.Sprintf("Auth error: %v",err))
			return
		}
	
		user,err := apiCfg.DB.GetUsersByApiKey(r.Context(),apiKey)
		if err != nil {
			responseWithError(w,400,fmt.Sprintf("Couldnt get user: %v",err))
			return
		}

		handler(w,r,user)
	}
}