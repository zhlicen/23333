package beepermission

type PermissionChecker interface {
	Check(params ...interface{}) error
}
