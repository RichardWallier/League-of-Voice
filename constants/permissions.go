package constants

// Richard Wallier: Based on db records, take care to update if the db records are updated
const (
	// Users
	PermissionUsersReadSelf   = "users.read.self"
	PermissionUsersReadAny    = "users.read.any"
	PermissionUsersCreate     = "users.create"
	PermissionUsersUpdateSelf = "users.update.self"
	PermissionUsersUpdateAny  = "users.update.any"
	PermissionUsersDeleteSelf = "users.delete.self"
	PermissionUsersDeleteAny  = "users.delete.any"

	// Roles
	PermissionRolesRead   = "roles.read"
	PermissionRolesWrite  = "roles.write"
	PermissionRolesDelete = "roles.delete"

	// Permissions
	PermissionPermissionsRead   = "permissions.read"
	PermissionPermissionsWrite  = "permissions.write"
	PermissionPermissionsDelete = "permissions.delete"
)
