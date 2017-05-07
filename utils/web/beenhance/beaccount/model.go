package beaccount

type AccountModel interface {
	Add(accountInfo *AccountInfo) error
	Delete(uid string) error

	GetUserId(idName IdName, loginId string) (*UserId, error)
	GetAccountInfo(uid string) (*AccountInfo, error)

	GetAccountBaseInfo(uid string) (*AccountBaseInfo, error)
	GetOAuth2Id(uid string) (map[KeyName]string, error)
	GetProfiles(uid string) (map[KeyName]string, error)
	GetOthers(uid string) (map[KeyName]string, error)
	GetAccountStatus(uid string) (*AccountStatus, error)

	UpdateAccountBaseInfo(uid string, baseInfo *AccountBaseInfo) error
	UpdateOAuth2Id(uid string, ids map[KeyName]string) error
	UpdateProfiles(uid string, profiles map[KeyName]string) error
	UpdateOthers(uid string, others map[KeyName]string) error
	UpdateAccountStatus(uid string, status *AccountStatus) error
}
