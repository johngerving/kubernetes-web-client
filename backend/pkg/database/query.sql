-- name: ListUsers :many
SELECT * FROM users;

-- name: FindUserWithEmail :one
SELECT * FROM users WHERE email = $1; 

-- name: CreateUser :exec
INSERT INTO users (email) VALUES ($1);

-- name: CreateWorkspace :exec
INSERT INTO workspaces (name, owner) VALUES ($1, $2);