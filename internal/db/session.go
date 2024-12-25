package db

import "time"

// Session holds the session data.
type Session struct {
	ID        string    `json:"id" bson:"_id"`
	CSRF      CSRF      `json:"csrf" bson:"csrf"`
	ExpiresAt time.Time `json:"expiresAt" bson:"expiresAt"`
}

// CSRF holds the CSRF token and its expiration time.
type CSRF struct {
	Token     string    `json:"token" bson:"token"`
	ExpiresAt time.Time `json:"expiresAt" bson:"expiresAt"`
}
