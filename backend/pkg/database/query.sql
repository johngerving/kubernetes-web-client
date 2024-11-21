-- name: ListUsers :many
SELECT * FROM users;

-- name: FindUserWithId :one
SELECT * FROM users WHERE id = $1;

-- name: FindUserWithEmail :one
SELECT * FROM users WHERE email = $1; 

-- name: CreateUser :exec
INSERT INTO users (email) VALUES ($1);

-- name: CreateWorkspace :one
INSERT INTO workspaces (name, owner) VALUES ($1, $2) RETURNING *;

-- name: ListUserWorkspaces :many
SELECT * FROM workspaces WHERE owner = $1;