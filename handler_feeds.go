package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/witczakxd/first-web-api/internal/database"
)


func (apiCfg *apiConfig)handlerCreateFeed(w http.ResponseWriter,r *http.Request,user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params :=parameters{}
	err := decoder.Decode(&params)
	if err!=nil {
		responseWithError(w,400,fmt.Sprintf("Error parsing json %v", err))
		return
	}

	feed,err := apiCfg.DB.CreateFeed(r.Context(),database.CreateFeedParams{
		ID:			uuid.New(),
		CreatedAt: 	time.Now().UTC(),
		UpdatedAt: 	time.Now().UTC(),
		Name: 		params.Name,
		Url: 		params.URL,
		UserID: 	user.ID,
	})

	if err!=nil {
		responseWithError(w,400,fmt.Sprintf("Couldnt create user %v", err))
		return
	}
	
	responseWithJSON(w,201,databaseFeedToFeed(feed))
}

func (apiCfg *apiConfig)handlerGetFeeds(w http.ResponseWriter,r *http.Request) {
	feeds,err := apiCfg.DB.GetFeeds(r.Context())

	if err!=nil {
		responseWithError(w,400,fmt.Sprintf("Couldnt create user %v", err))
		return
	}
	
	responseWithJSON(w,201,databaseFeedsToFeeds(feeds))
}
