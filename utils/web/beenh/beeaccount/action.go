package beeaccount

// AccountAction account action
type AccountAction int

// Account Actions definition
const (
	AccountAction_Begin AccountAction = iota
	ActionAccountLogoutSession
	ActionAccountChangePwd
	ActionAccountGetAccountBasicInfo
	ActionAccountUpdateAccountBasicInfo
	ActionAccountGetProfiles
	ActionAccountUpdateProfiles
	ActionAccountGetOthers
	ActionAccountUpdateOthers
	AccountActionEnd
)
