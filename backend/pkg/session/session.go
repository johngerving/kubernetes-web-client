package session

import (
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewStore takes a pointer to a pgxpool.Pool struct instance, configures a session manager
// with pgsxstore as the session store, and returns a pointer to the session manager.
func NewStore(pool *pgxpool.Pool) *scs.SessionManager {
	// Initialize a new session manager and configure it to use pgxstore as the session store
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(pool)
	sessionManager.Lifetime = 24 * time.Hour // Set session lifetime to 1 day

	return sessionManager
}
