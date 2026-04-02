package psql

import (
	"context"
	"fmt"

	"github.com/0bby/genki/internal/models"
)

func (p *Psql) UpsertSyncCursor(ctx context.Context, c *models.SyncCursor) error {
	_, err := p.conn.Exec(ctx, `
		INSERT INTO sync_cursors (user_id, source, resource, last_synced_at, cursor_value)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, source, resource) DO UPDATE
		SET last_synced_at = EXCLUDED.last_synced_at, cursor_value = EXCLUDED.cursor_value`,
		c.UserID, c.Source, c.Resource, c.LastSyncedAt, c.CursorValue)
	return err
}

func (p *Psql) GetSyncCursor(ctx context.Context, userID int32, source, resource string) (*models.SyncCursor, error) {
	row := p.conn.QueryRow(ctx, `
		SELECT id, user_id, source, resource, last_synced_at, cursor_value
		FROM sync_cursors WHERE user_id = $1 AND source = $2 AND resource = $3`, userID, source, resource)
	c := &models.SyncCursor{}
	err := row.Scan(&c.ID, &c.UserID, &c.Source, &c.Resource, &c.LastSyncedAt, &c.CursorValue)
	if err != nil {
		return nil, fmt.Errorf("GetSyncCursor: %w", err)
	}
	return c, nil
}

func (p *Psql) GetAllSyncCursors(ctx context.Context, userID int32) ([]models.SyncCursor, error) {
	rows, err := p.conn.Query(ctx, `
		SELECT id, user_id, source, resource, last_synced_at, cursor_value
		FROM sync_cursors WHERE user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("GetAllSyncCursors: %w", err)
	}
	defer rows.Close()
	var result []models.SyncCursor
	for rows.Next() {
		var c models.SyncCursor
		if err := rows.Scan(&c.ID, &c.UserID, &c.Source, &c.Resource, &c.LastSyncedAt, &c.CursorValue); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, nil
}

func (p *Psql) UpsertOAuthToken(ctx context.Context, t *models.OAuthToken) error {
	_, err := p.conn.Exec(ctx, `
		INSERT INTO oauth_tokens (user_id, provider, access_token, refresh_token, expires_at, scopes)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, provider) DO UPDATE
		SET access_token = EXCLUDED.access_token, refresh_token = EXCLUDED.refresh_token,
		    expires_at = EXCLUDED.expires_at, scopes = EXCLUDED.scopes`,
		t.UserID, t.Provider, t.AccessToken, t.RefreshToken, t.ExpiresAt, t.Scopes)
	return err
}

func (p *Psql) GetOAuthToken(ctx context.Context, userID int32, provider string) (*models.OAuthToken, error) {
	row := p.conn.QueryRow(ctx, `
		SELECT id, user_id, provider, access_token, refresh_token, expires_at, scopes
		FROM oauth_tokens WHERE user_id = $1 AND provider = $2`, userID, provider)
	t := &models.OAuthToken{}
	err := row.Scan(&t.ID, &t.UserID, &t.Provider, &t.AccessToken, &t.RefreshToken, &t.ExpiresAt, &t.Scopes)
	if err != nil {
		return nil, fmt.Errorf("GetOAuthToken: %w", err)
	}
	return t, nil
}

func (p *Psql) DeleteOAuthToken(ctx context.Context, userID int32, provider string) error {
	_, err := p.conn.Exec(ctx, `DELETE FROM oauth_tokens WHERE user_id = $1 AND provider = $2`, userID, provider)
	return err
}
