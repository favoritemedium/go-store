package store

import (
  "fmt"
	"time"
	"errors"
	"strings"
	"database/sql"
	"github.com/go-sql-driver/mysql"
)

var (
	ErrEmailRequired = errors.New("email field may not be blank")
	ErrFullNameRequired = errors.New("full name field may not be blank")
	ErrNameToUseRequired = errors.New("name to use filed may not be blank")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrInvalidEmailVerifyCode = errors.New("invalid email verification code")
	ErrDuplicateEmail = errors.New("email already registered")
)

// Type User is the basic user type.
type User struct {
	Id uint32         `json:"id"`
	Email string      `json:"email"`
	FullName string   `json:"fullName"`
	NameToUse string  `json:"nameToUse"`
	IsActive bool     `json:"isActive"`
	IsAdmin bool      `json:"isAdmin"`
	IsSuperUser bool  `json:"isSuperUser"`
	CreatedAt int64   `json:"createdAt"`
	UpdatedAt int64   `json:"updatedAt"`

  allowedauth uint8
	pwhash string
}

// UserCreate makes a new user record.  AuthUser must either have admin
// privileges or have the special NewUser privilege.
func UserCreate(creator AuthUser, u *User) error {
	if !(creator.IsAdmin() || creator.IsSuperUser() || creator.IsNewUser()) {
		return ErrUnauthorized
	}

	if !creator.IsAdmin() { u.IsAdmin = false }
	if !creator.IsSuperUser() { u.IsSuperUser = false }

	if creator.IsNewUser() {
		u.Email = creator.GetEmail()
		u.allowedauth = creator.GetProvider()
	}

	if err := u.Validate(); err != nil {
		return err
	}

	u.CreatedAt = time.Now().Unix()
	u.UpdatedAt = u.CreatedAt

	query :=	"INSERT INTO " + UserTable + " (email, allowedauth, pwhash, fullname, nametouse, isactive, isadmin, issuperuser, createdat, updatedat) VALUES (?,?,?,?,?,?,?,?,?,?)"
	r, err := db.Exec(query, u.Email, u.allowedauth, u.pwhash, u.FullName, u.NameToUse, u.IsActive, u.IsAdmin, u.IsSuperUser, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		if err0, ok := err.(*mysql.MySQLError); ok {
			// 1062 is mysql for unique contstraint violation
			if err0.Number == 1062 {
				return ErrDuplicateEmail
			}
		}
		return err
	}

  x, _ := r.LastInsertId()
  u.Id = uint32(x)
  return nil
}

const userSelect = "SELECT id, email, fullname, nametouse, isactive, isadmin, issuperuser, createdat, updatedat FROM " + UserTable

func readRows(rows *sql.Rows, err error) StoreChannel {
	ch := make(StoreChannel)
	go func() {
		defer close(ch)
		if err != nil {
			ch <- StoreResult{User{}, err}
			return
		}
		defer rows.Close()
		for rows.Next() {
		  var u User
			err := rows.Scan(&u.Id, &u.Email, &u.FullName, &u.NameToUse, &u.IsActive,
				&u.IsAdmin, &u.IsSuperUser, &u.CreatedAt, &u.UpdatedAt)
			ch <- StoreResult{u, err}
		}
		if err := rows.Err(); err != nil {
			ch <- StoreResult{User{}, err}
		}
	}()
	return ch
}

// UserRead gets one user record by primary key.
func UserRead(id uint32) StoreChannel {
	query := userSelect + " WHERE id=? LIMIT 1"
	return readRows(db.Query(query, id))
}

// UserReadn gets one range of user records by primary key.
// Returns up to n records starting from id=start.
func UserReadn(start, n uint32) StoreChannel {
	query := fmt.Sprintf(userSelect + " WHERE id>=? ORDER BY id LIMIT %d", n)
	return readRows(db.Query(query, start))
}

// UserReadMultiple gets an arbitrary set of records by primary key.
func UserReadMultiple(ids []uint32) StoreChannel {
	if len(ids) == 0 {
		return emptyStoreChannel()
	}
	strids := make([]string, len(ids))
	for _, id := range ids {
		strids = append(strids, fmt.Sprintf("%d", id))
	}
	idlist := strings.Join(strids, ",")
	query := userSelect + " WHERE id IN (" + idlist + ") ORDER BY id"
	return readRows(db.Query(query))
}

// UserSearch finds all records in which any of the fields match the search text.
// Partial matches may be achieved using the % wildcard, e.g. "wee%" will match
// any value that begins with "wee".  Results are ordered by activity - recently
// active users are shown first
func UserSearch(searchtext string, fields []string) StoreChannel {
	// TODO validate fields
	if len(fields) == 0 || len(searchtext) == 0 {
		return emptyStoreChannel() // TODO make this an error
	}
	effs := make([]string, len(fields))
	texts := make([]interface{}, len(fields))
	for _, field := range fields {
		effs = append(effs, field + " LIKE ?")
		texts = append(texts, searchtext)
	}
	// TODO fix ordering
	query := userSelect + " WHERE " + strings.Join(effs, " OR ") + " ORDER BY " + fields[0]
	return readRows(db.Query(query, texts...))
}

// UserSearchn is the same as Search, except that a subset of results is returned:
// skip is the number of initial records to skip, and n is the number of records
// to return.
func UserSearchn(searchtext string, fields []string, skip, n uint32) StoreChannel {
	return emptyStoreChannel()
}

// SetPassword sets the password on this user.
// Call Update after this to write the new password hash to the db.
func (u *User) SetPassword(newpw string) error {
	return ErrNotImplemented
}

// SetEmail sets the email address for this user.  This method must be used
// rather than setting the value directly when changing one's own email address
// (i.e. User is the same as AuthUser).
func (u *User) SetEmail(newemail, verifyCode string) error {
	return ErrNotImplemented
}

// Update saves the user record.  Must supply a valid AuthUser with the
// requisite permissions.
//
// Fields Email, FullName, NameToUse, IsActive may be updated by the owner or
// by an account with higher privileges, i.e. admin or superuser for regular
// accounts, and only superuser for admin accounts.
//
// Adding "Password" as one of the fields will cause the password set by
// SetPassword() to be the new password for the user.
//
// Field IsAdmin may be changed to true by any admin and to false only by
// a superuser.
//
// Field IsSuperUser may only be changed by a superuser.
// The first user (id=1) is always a superuser and can never be demoted.
//
// Fields CreatedAt and UpdatedAt are set automatically.
func (u *User) Update(updater AuthUser, fields []string) StoreResult {
	return StoreResult{}
}

// Delete permanently deletes the user record.  Use with care.  Must supply a
// valid AuthUser that owns the record or has higher privileges than the
// record owner.
func (u *User) Delete(deleter AuthUser) error {
	return ErrNotImplemented
}

func (u *User) Validate() error {
	if u.Email == "" {
		return ErrEmailRequired
	}
	if u.FullName == "" {
		return ErrFullNameRequired
	}
	if u.NameToUse == "" {
		return ErrNameToUseRequired
	}
	return nil
}
