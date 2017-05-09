package beeaccount

// AccountModel model interface for account manager
type AccountModel interface {
	Add(accountInfo *AccountInfo) error
	Delete(uid string) error

	GetUserId(idName string, loginId string) (*UserId, error)
	GetAccountInfo(uid string) (*AccountInfo, error)

	GetAccountBasicInfo(uid string) (*AccountBasicInfo, error)
	GetOAuth2Id(uid string) (map[string]string, error)
	GetProfiles(uid string) (map[string]string, error)
	GetOthers(uid string) (map[string]string, error)
	GetAccountStatus(uid string) (*AccountStatus, error)

	UpdateAccountBasicInfo(uid string, basicInfo *AccountBasicInfo) error
	UpdateOAuth2Id(uid string, ids map[string]string) error
	UpdateProfiles(uid string, profiles map[string]string) error
	UpdateOthers(uid string, others map[string]string) error
	UpdateAccountStatus(uid string, status *AccountStatus) error
}
