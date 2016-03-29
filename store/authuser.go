package store

import (
	"net/http"
	"errors"
)

const (
	// AuthTokenMaxIdleHours is the maximum length of time an auth token can
	// remain valid without any activity.
	AuthTokenMaxIdleHours = 24

	// RefreshTokenMaxAgeDays is the how far in the future refresh tokens are
	// set to expire.  Note that refresh tokens are invalidated immediately
	// when they're used.
	RefreshTokenMaxAgeDays = 30

	// EmailVerificationMaxAgeHours is how much an email verification token is valid for.
	// This is how much time a user has to respond to their welcom email.
	EmailVerificationMaxAgeHours = 24

	// SigninEventMaxAgeDays is the length of time signin history should be kept.
	SigninEventMaxAgeDays = 90
)

var (
	ErrInvalidEmailPassword = errors.New("email and password don't match")
	ErrInvalidGoogleIdToken = errors.New("invalid Google ID token")
	ErrInvalidFacebookAccessToken = errors.New("invalid Facebook access token")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
	ErrInvalidAuthToken = errors.New("invalid or expired auth token")
)

// RegisterClientID saves the Client ID for an authentication provider such as
// Google for use in handling that provider's sign-ins.
// Must be called for each provider before that provider is used.
func RegisterOauthClientId(provider int, clientID string) {
}

// Type AuthUser represents the current signed-in user.
type AuthUser struct {
	FullName string
	NameToUse string

	id uint32
	provider uint8
	email string
	authToken string
	authTokenExpiry int64
	refreshToken string
	refreshTokenExpiry int64
	isAdmin bool
	isSuperUser bool
	isNewUser bool
}

// Type Activity represents the connection details of one http request.
// It's derived from the headers on the request object.
type Activity struct{
	Time uint64
	IP string
	Device string  // e.g. "Firefox on Linux"
}

// ActivityFromRequest takes an http.Request and creates an Activity based on the
// the IP and header information contained in it.
func ActivityFromRequest(r *http.Request) *Activity {
	return &Activity{}
}

// SigninEmail signs in a user via email and password.
func SigninEmail(email, password string, a Activity) (*AuthUser, error) {
	return nil, ErrNotImplemented
}

// SigninOauth signs in a user using the id/access token from an oauth2 provider
// such as Google or Facebook.
func SigninOauth(provider uint8, token string, a Activity) (*AuthUser, error) {
	return nil, ErrNotImplemented
}

// SigninRefresh re-signs-in a previously signed-in user.
// If the Device field doesn't match, that suggest that the token was compromised.
// In such a case, the sign-in will fail and the token will be invalidated.
func SigninRefresh(refreshToken string, a Activity) (*AuthUser, error) {
	return nil, ErrNotImplemented
}

// GetEmailVerifyCode initiations a new user process by taking a user-supplied
// email address and generating a verification code that can be emailed to the
// new user, normally as part of a "click to verify" link.
func GetEmailVerifyCode(email string) string {
	return ""
}

// NewUserEmail takes a verification code from GetEmailVerifyCode and returns
// an AuthUser object with the special privilege of being allowed to create
// one User.
// An error is returned if this user already exists in the database.
func NewUserEmail(verifyCode string, a Activity) (*AuthUser, error) {
	return nil, ErrNotImplemented
}

// NewUserGoogle takes a token from an oauth2 provider such as Google and
// returns an AuthUser object with the special privilege of being allowed
// to create one User.
// An error is returned if this user already exists in the database.
func NewUserOauth(provider uint8, token string, a Activity) (*AuthUser, error) {
	return nil, ErrNotImplemented
}

// VerifySession connects an auth token to its user.
// This would normally be called at the beginning of every request.
// If the Device field doesn't match, that suggests that the token was compromised.
// In such a case, the user is automatically signed out, and the refresh token is
// invalidated, forcing the user to re-authenticate.
func VerifySession(authToken string, a Activity) (*AuthUser, error) {
	return nil, ErrNotImplemented
}

// GetId returns the internal ID of this user; useful when it's a foreign key
// elsewhere.
func (u *AuthUser) GetId() uint32 {
	return u.id
}

// GetProvider returns the authentication provider for this user.
func (u *AuthUser) GetProvider() uint8 {
	return u.provider
}

// GetEmail returns the verified email address of this user.
func (u *AuthUser) GetEmail() string {
	return u.email
}

// IsAdmin tells us if this user has admin privileges.
func (u *AuthUser) IsAdmin() bool {
	return u.isAdmin
}

// IsSuperUser tells us if this user has superuser privileges.
func (u *AuthUser) IsSuperUser() bool {
	return u.isSuperUser || u.id == 1
}

// IsNewUser tells us if there is not yet a matching User record, and we
// can create one.
func (u *AuthUser) IsNewUser() bool {
	return u.isNewUser
}

// GetAuthToken retrieves a token that can be used later with VerifySession.
func (u *AuthUser) GetAuthToken() string {
	return u.authToken
}

// GetRefreshToken retrieves a token that can be used later with SigninRefresh.
func (u *AuthUser) GetRefreshToken() string {
	return u.refreshToken
}

// GetRefreshTokenExpiry retrieves the expiration time (unix time) of the refresh
// token. It's not necessary for the caller to keep track of this, as an expired
// refresh token will simply get rejected, but it's useful to set the expiration
// time of cookies, for example.
func (u *AuthUser) GetRefreshTokenExpiry() int64 {
	return u.refreshTokenExpiry
}

// GetThisSignin returns the time, IP and device of the current signin.
func (u *AuthUser) GetThisSignin() *Activity {
	return nil
}

// GetPreviousSignins finds up to n recent signin events and returns their details.
// Signin events are automatically deleted after some time and may not be available
func (u *AuthUser) GetPreviousSignins(n int) []*Activity {
	return []*Activity{}
}
