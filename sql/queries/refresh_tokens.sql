-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id)
VALUES (
	$1,
	$2
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * 
FROM refresh_tokens
WHERE token = $1;

-- name: UpdateRefreshToken :exec
UPDATE refresh_tokens
SET
    token = $2,
    updated_at = NOW(),
    expires_at = NOW() + INTERVAL '60 days'
WHERE token = $1;