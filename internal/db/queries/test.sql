-- name: GetUser :one
SELECT * FROM test_users
WHERE email = $1;