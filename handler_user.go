package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/OmarEP/blog-aggregator/internal/database"
	"github.com/google/uuid"
)



func (apiCfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name 	string 	`json:"name"`
	}
	
	type response struct {
		ID			string		`json:"id"`
		Created_At 	time.Time	`json:"created_at"`
		Updated_At 	time.Time	`json:"updated_at"`
		Name 		string		`json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error passing JSON: %v", err))
		return 
	}
	
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:	uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't create user: %v", err))
		return 
	}

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request, user database.User){
	// apiKey, err := auth.GetAPIKey(r.Header)
	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, "Couldn't find api key")
	// 	return 
	// }

	// user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
	// if err != nil {
	// 	respondWithError(w, http.StatusNotFound, "Couldn't get user")
	// 	return 
	// }

	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}