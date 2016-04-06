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
	"github.com/go-sql-driver/mysql"
	"os"
	"fmt"
	"errors"
)

// Names of tables in the database that are used by this package.
// This package does not concern itself with creating these tables;
// that must be done elsewhere.
const (
	UserTable = "users"
	UserSigninsTable = "users_signins"
	UserTokensTable = "users_tokens"
	UserVerifyTable = "users_verify"

)

// These are the possible user roles.
// Users may have any combination of roles.
const (
	SuperRole = 1 << iota
	AdminRole
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

// Type StoreResult contains one row of returned data from db queries
type StoreResult struct {
	Data interface{}
	Err  error
}

// Type StoreChannel is used to return a set of data from db queries
type StoreChannel chan StoreResult


var db *sql.DB

// Init sets the db handle for use in all subsequent db operations.
// Must be called once.
func Init(dbhandle *sql.DB) {
	db = dbhandle
}

// Catch aborts on an error
func Catch(err error) {
	if err == nil { return; }
	fmt.Println(err)
	os.Exit(1)
}

// IsDuplicate tests err to see if it's a unique constraint db error
func IsDuplicate(err error) bool {
	//if err != nil {
		if err0, ok := err.(*mysql.MySQLError); ok {
			// 1062 is mysql for unique contstraint violation
			return err0.Number == 1062
		}
	//}
	return false
}

// a placeholder to simulate no results
func emptyStoreChannel() StoreChannel {
	ch := make(StoreChannel)
  go close(ch)
  return ch
}
