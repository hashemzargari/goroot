package goroot

type ApiHttpMethod string

func (m ApiHttpMethod) String() string {
	return string(m)
}

const (
	Get    ApiHttpMethod = "GET"
	Post   ApiHttpMethod = "POST"
	Patch  ApiHttpMethod = "PATCH"
	Put    ApiHttpMethod = "PUT"
	Delete ApiHttpMethod = "DELETE"
)

type Handler interface {
	// Handle is the main method of the handler
	Handle(ctx Context, request any) (response any, err error)

	// GetRequestType returns the type of the request
	GetRequestType() any

	// GetResponseType returns the type of the response
	GetResponseType() any

	// GetPermissionClaims returns the permission claims required to execute the handler
	GetPermissionClaims() []PermissionClaim

	GetAuthenticationClaims() AuthenticationClaim

	// GetApiRoute returns the api route for the handler
	GetApiRoute() string

	// GetMethod returns the method for the handler
	GetMethod() ApiHttpMethod
}
