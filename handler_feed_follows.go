package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/witczakxd/first-web-api/internal/database"
)


func (apiCfg *apiConfig)handlerCreateFeedFollow(w http.ResponseWriter,r *http.Request,user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params :=parameters{}
	err := decoder.Decode(&params)
	if err!=nil {
		responseWithError(w,400,fmt.Sprintf("Error parsing json %v", err))
		return
	}

	feedFollow,err := apiCfg.DB.CreateFeedFollow(r.Context(),database.CreateFeedFollowParams{
		ID:			uuid.New(),
		CreatedAt: 	time.Now().UTC(),
		UpdatedAt: 	time.Now().UTC(),
		UserID: 	user.ID,
		FeedID: 	params.FeedID,
	})

	if err!=nil {
		responseWithError(w,400,fmt.Sprintf("Couldnt create feed_follow %v", err))
		return
	}
	
	responseWithJSON(w,201,databaseFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *apiConfig)handlerGetFeedFollow(w http.ResponseWriter,r *http.Request,user database.User) {

	feedFollows,err := apiCfg.DB.GetFeedFollows(r.Context(),user.ID)
	if err!=nil {
		responseWithError(w,400,fmt.Sprintf("Couldnt get feedfollows %v", err))
		return
	}
	
	responseWithJSON(w,201,databaseFeedFollowsToFeedFollows(feedFollows))
}

func (apiCfg *apiConfig)handlerDeleteFeedFollow(w http.ResponseWriter,r *http.Request,user database.User) {

	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID,err := uuid.Parse(feedFollowIDStr);
	if err != nil {
		responseWithError(w,400,fmt.Sprintf("Couldnt parse feedFollowID %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(),database.DeleteFeedFollowParams{
		ID: feedFollowID,
		UserID: user.ID,
	})
	fmt.Println(feedFollowID.String())
	fmt.Println(user.ID.String())
	if err!=nil {
		responseWithError(w,400,fmt.Sprintf("Couldnt delete feedfollow %v", err))
		return
	}

	responseWithJSON(w,200,struct{}{})
}

