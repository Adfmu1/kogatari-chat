-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id)
VALUES (
	$1,
	$2
)
RETURNING *;