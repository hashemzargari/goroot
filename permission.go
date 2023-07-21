package goroot

type PermissionClaim string

type PermissionBackend interface {
	// HasPermissions returns true if the user has the required permissions
	HasPermissions(claims []PermissionClaim) bool

	// GetPermissionClaims returns the permission claims for the user
	GetPermissionClaims() []PermissionClaim
}
