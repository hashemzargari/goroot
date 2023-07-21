package goroot

type User struct {
	ID       uint
	Username string
}

type AuthenticationClaim int

const (
	AllowAny AuthenticationClaim = 1 << iota
	AllowAuthenticated
	AllowAdmin
	AllowOwner
)

type AuthenticationBackend interface {
	// Authenticate authenticates the user
	Authenticate(ctx Context, claims AuthenticationClaim) (user User, isAllowed bool, err error)
}
