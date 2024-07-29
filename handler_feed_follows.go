package main 

import(
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/OmarEP/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig)handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID 	uuid.UUID 	`json:"feed_id"`
	}
	
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error passing JSON: %v", err))
		return 
	}
	
	feed_follow, err := apiCfg.DB.CreateFeedsFollow(r.Context(), database.CreateFeedsFollowParams{
		ID:	uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: params.FeedID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't create feed follow: %v", err))
		return 
	}

	respondWithJSON(w, http.StatusOK, databaseFeedFollowToFeedFollow(feed_follow))
}

func (apiCfg *apiConfig) handlerFeedFollowsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	feed_follows, err := apiCfg.DB.GetFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return 
	}

	respondWithJSON(w, http.StatusOK, databaseFeedFollowsToFeedFollows(feed_follows))
}

func (apiCfg *apiConfig) handlerFeedFollowDelete(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollowIDStr := r.PathValue("feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid feed follow ID")
		return 
	}
	
	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		ID: feedFollowID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete feed follow")
		return 
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}