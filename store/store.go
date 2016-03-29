// Package store provides the db interface and encapsulates all the high-level
// operations that need to be done on the db.
//
// The package includes methods to hand user validation, including creating
// user accounts (SSO or email+password), validating returning users, and
// managing sessions.  The AuthUser type represents an authenticated user.
//
// For db operations that are restricted, including all update operations,
// a valud AuthUser with the required privileges must be supplied.
// That way, all methods here can be exposed directly to an external API.
//
// For better perfomance, we may want to keep the session data in a cache.
// This is an implementation detail and will not change the usage.
package store

import (
	"database/sql"
	"errors"
)

// Names of tables in the database that are used by this package.
// This package does not concern itself with creating these tables;
// that must be done elsewhere.
const (
	UserTable = "users"
	SignedinUserTable = "users_signedin"
	SigninHistoryTable = "users_history"
	EmailVerifyTable = "users_verify_email"

)

// These are the authentication providers we know about.
// EmailAuth means no provider, just email & password.
const (
	EmailAuth = 1 << iota
	GoogleAuth
	FacebookAuth
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

// remove this
var ErrNotImplemented = errors.New("not implemented")

// Type StoreResult contains returned data from database queries
type StoreResult struct {
	Data interface{}
	Err  error
}

// Type StoreError is used for any db-level erros
type StoreError struct {
	message string
}

func (e *StoreError) Error() string {
	return e.message
}

type StoreChannel chan StoreResult

var db *sql.DB

// Init sets the db handle for use in all subsequent db operations.
// Must be called once.
func Init(dbhandle *sql.DB) {
	db = dbhandle
}

func emptyStoreChannel() StoreChannel {
	ch := make(StoreChannel)
  go close(ch)
  return ch
}
