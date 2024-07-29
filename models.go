package main

import (
	"time"
	"database/sql"

	"github.com/OmarEP/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

type Feed struct {
	ID        		uuid.UUID `json:"id"`
	CreatedAt 		time.Time `json:"created_at"`
	UpdatedAt 		time.Time `json:"updated_at"`
	Name      		string    `json:"name"`
	Url       		string    `json:"url"`
	UserID    		uuid.UUID `json:"user_id"`
	LastFetchedAt 	*time.Time `json:"last_fetched_at"`
}

type FeedFollow struct {
	ID        uuid.UUID	`json:"id"`
	CreatedAt time.Time	`json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID	`json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

type Post struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Title       string     `json:"title"`
	Url         string     `json:"url"`
	Description *string    `json:"description"`
	PublishedAt *time.Time `json:"published_at"`
	FeedID      uuid.UUID  `json:"feed_id"`
}

func databaseUserToUser(user database.User) User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	}
}

func databaseFeedToFeed(feed database.Feed) Feed {
	return Feed{
		ID:        		feed.ID,
		CreatedAt: 		feed.CreatedAt,
		UpdatedAt: 		feed.UpdatedAt,
		Name:      		feed.Name,
		Url:       		feed.Url,
		UserID:    		feed.UserID,
		LastFetchedAt:	nullTimeToTimePtr(feed.LastFetchedAt),
	}
}

func nullTimeToTimePtr(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func nullStringToStringPtr(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}


func databaseFeedsToFeeds(feeds []database.Feed) []Feed {
	response := make([]Feed, len(feeds))
	for i, feed := range(feeds) {
		response[i] = databaseFeedToFeed(feed)
	}
	return response
}

func databaseFeedFollowToFeedFollow(feed_follow database.FeedFollow) FeedFollow{
	return FeedFollow{
		ID: feed_follow.ID,
		CreatedAt: feed_follow.CreatedAt,
		UpdatedAt: feed_follow.UpdatedAt,
		UserID: feed_follow.UserID,
		FeedID: feed_follow.FeedID,
	}
}

func databaseFeedFollowsToFeedFollows(feed_follows []database.FeedFollow) []FeedFollow {
	response := make([]FeedFollow, len(feed_follows))
	for i, feed_follow := range(feed_follows) {
		response[i] = databaseFeedFollowToFeedFollow(feed_follow)
	}
	return response
}

func databasePostToPost(post database.Post) Post {
	return Post{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Url:         post.Url,
		Description: nullStringToStringPtr(post.Description),
		PublishedAt: nullTimeToTimePtr(post.PublishedAt),
		FeedID:      post.FeedID,
	}
}

func databasePostsToPosts(posts []database.Post) []Post {
	result := make([]Post, len(posts))
	for i, post := range posts {
		result[i] = databasePostToPost(post)
	}
	return result
}
