package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/7-Dany/dev/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	params := parameters{}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("failed to create a feed follow: %v", err))
		return
	}

	respondWithJSON(w, 201, dbFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't get feeds, %v", err))
		return
	}
	respondWithJSON(w, 200, dbFeedFollowsToFeedFollows(feedFollows))
}

func (apiCfg *apiConfig) handleDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIdStr := chi.URLParam(r, "feedFollowID")
	feedFollowId, err := uuid.Parse(feedFollowIdStr)
	if err != nil {
		respondWithJSON(w, 400, fmt.Sprintf("couldn't parse feed id, %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowId,
		UserID: user.ID,
	})
	if err != nil {
		respondWithJSON(w, 400, fmt.Sprintf("error deleting feed follow, %v", err))
		return
	}

	respondWithJSON(w, 200, struct{}{})
}
