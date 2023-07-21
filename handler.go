package goroot

type ApiHttpMethod string

const (
	Get    ApiHttpMethod = "GET"
	Post   ApiHttpMethod = "POST"
	Patch  ApiHttpMethod = "PATCH"
	Put    ApiHttpMethod = "PUT"
	Delete ApiHttpMethod = "DELETE"
)

type Handler interface {
	// Handle is the main method of the handler
	Handle(ctx Context, req interface{}) (interface{}, error)

	// GetRequestType returns the type of the request
	GetRequestType() interface{}

	// GetResponseType returns the type of the response
	GetResponseType() interface{}

	// GetPermissionClaims returns the permission claims required to execute the handler
	GetPermissionClaims() []PermissionClaim

	GetAuthenticationClaims() AuthenticationClaim

	// GetApiRoute returns the api route for the handler
	GetApiRoute() string

	// GetMethod returns the method for the handler
	GetMethod() ApiHttpMethod
}
