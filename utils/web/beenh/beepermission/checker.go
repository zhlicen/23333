package beepermission

// PermissionChecker permission check interface
type PermissionChecker interface {
	Check(action interface{}, params ...interface{}) error
}
