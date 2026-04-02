-- name: UpsertSyncCursor :exec
INSERT INTO sync_cursors (user_id, source, resource, last_synced_at, cursor_value)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id, source, resource) DO UPDATE
SET last_synced_at = EXCLUDED.last_synced_at, cursor_value = EXCLUDED.cursor_value;

-- name: GetSyncCursor :one
SELECT * FROM sync_cursors
WHERE user_id = $1 AND source = $2 AND resource = $3;

-- name: GetAllSyncCursors :many
SELECT * FROM sync_cursors WHERE user_id = $1;

-- name: UpsertOAuthToken :exec
INSERT INTO oauth_tokens (user_id, provider, access_token, refresh_token, expires_at, scopes)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (user_id, provider) DO UPDATE
SET access_token = EXCLUDED.access_token,
    refresh_token = EXCLUDED.refresh_token,
    expires_at = EXCLUDED.expires_at,
    scopes = EXCLUDED.scopes;

-- name: GetOAuthToken :one
SELECT * FROM oauth_tokens
WHERE user_id = $1 AND provider = $2;

-- name: DeleteOAuthToken :exec
DELETE FROM oauth_tokens WHERE user_id = $1 AND provider = $2;
