package beeaccount

// AccountModel model interface for account manager
type AccountModel interface {
	Add(accountInfo *AccountInfo) error
	Delete(uid string) error

	GetUserId(idName IdName, loginId string) (*UserId, error)
	GetAccountInfo(uid string) (*AccountInfo, error)

	GetAccountBasicInfo(uid string) (*AccountBasicInfo, error)
	GetOAuth2Id(uid string) (map[KeyName]string, error)
	GetProfiles(uid string) (map[KeyName]string, error)
	GetOthers(uid string) (map[KeyName]string, error)
	GetAccountStatus(uid string) (*AccountStatus, error)

	UpdateAccountBasicInfo(uid string, basicInfo *AccountBasicInfo) error
	UpdateOAuth2Id(uid string, ids map[KeyName]string) error
	UpdateProfiles(uid string, profiles map[KeyName]string) error
	UpdateOthers(uid string, others map[KeyName]string) error
	UpdateAccountStatus(uid string, status *AccountStatus) error
}
