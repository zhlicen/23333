package beeaccount

type ModelQuery int

// AccountModel model interface for account manager
type AccountModel interface {
	Add(accountInfo *AccountInfo) error
	Delete(uid string) error

	GetUserID(idName string, loginID string) (*UserID, error)
	GetAccountInfo(uid string) (*AccountInfo, error)

	GetAccountBasicInfo(uid string) (*AccountBasicInfo, error)
	GetOAuth2ID(uid string) (map[string]string, error)
	GetProfiles(uid string) (map[string]string, error)
	GetOthers(uid string) (map[string]string, error)
	GetAccountStatus(uid string) (*AccountStatus, error)

	UpdateAccountBasicInfo(uid string, basicInfo *AccountBasicInfo) error
	UpdateOAuth2ID(uid string, ids map[string]string) error
	UpdateProfiles(uid string, profiles map[string]string) error
	UpdateOthers(uid string, others map[string]string) error
	UpdateAccountStatus(uid string, status *AccountStatus) error

	ModelQuery(id ModelQuery, params ...interface{}) (interface{}, error)
}
