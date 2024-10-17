-- name: ListUsers :many
SELECT * FROM users;

-- name: FindUserWithEmail :one
SELECT * FROM users WHERE email = $1; 

-- name: CreateUser :exec
INSERT INTO users (email) VALUES ($1);