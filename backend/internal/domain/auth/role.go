package auth

type Role string

const (
	RoleOwner  Role = "OWNER"
	RoleViewer Role = "VIEWER"
)
