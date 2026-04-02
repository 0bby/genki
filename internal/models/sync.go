package models

import "time"

type SyncCursor struct {
	ID           int32     `json:"id"`
	UserID       int32     `json:"user_id"`
	Source       string    `json:"source"`
	Resource     string    `json:"resource"`
	LastSyncedAt time.Time `json:"last_synced_at"`
	CursorValue  *string   `json:"cursor_value,omitempty"`
}

type OAuthToken struct {
	ID           int32     `json:"id"`
	UserID       int32     `json:"user_id"`
	Provider     string    `json:"provider"`
	AccessToken  string    `json:"-"`
	RefreshToken string    `json:"-"`
	ExpiresAt    time.Time `json:"expires_at"`
	Scopes       *string   `json:"scopes,omitempty"`
}
