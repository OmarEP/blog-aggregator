package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/OmarEP/blog-aggregator/internal/database"
	"github.com/google/uuid"
)



func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name 	string 	`json:"name"`
		Url		string	`json:"url"`
	}
	
	
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error passing JSON: %v", err))
		return 
	}
	
	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:	uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
		Url: params.Url,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't create user: %v", err))
		return 
	}

	feed_follow, err := apiCfg.DB.CreateFeedsFollow(r.Context(), database.CreateFeedsFollowParams{
		ID:	uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't create feed follow: %v", err))
		return 
	}

	respondWithJSON(w, http.StatusOK, struct {
		Feed			Feed 			`json:"feed"`
		FeedFollow		FeedFollow		`json:"feed_follow"`
	}{
		Feed: databaseFeedToFeed(feed),
		FeedFollow: databaseFeedFollowToFeedFollow(feed_follow),
	})
}

func (apiConfig *apiConfig) handlerFeedsGet(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiConfig.DB.GetFeeds(r.Context()) 
	if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't get feeds")
			return 
	}
	respondWithJSON(w, http.StatusOK, databaseFeedsToFeeds(feeds))
}