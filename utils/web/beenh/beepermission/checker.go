package beepermission

// PermissionChecker permission check interface
type PermissionChecker interface {
	Check(params ...interface{}) error
}
