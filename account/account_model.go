package account

type AccountModel interface {
	Add(accountInfo *AccountInfo) error
	Delete(uid AccountUid) error

	GetUidById(idName IdName, id string) (AccountUid, error)
	GetAccountInfo(uid AccountUid) (*AccountInfo, error)

	GetAccountBaseInfo(uid AccountUid) (*AccountBaseInfo, error)
	GetOAuth2Id(uid AccountUid) (map[KeyName]string, error)
	GetProfiles(uid AccountUid) (map[KeyName]string, error)
	GetOthers(uid AccountUid) (map[KeyName]string, error)
	GetAccountStatus(uid AccountUid) (*AccountStatus, error)

	UpdateAccountBaseInfo(uid AccountUid, baseInfo *AccountBaseInfo) error
	UpdateOAuth2Id(uid AccountUid, ids map[KeyName]string) error
	UpdateProfiles(uid AccountUid, profiles map[KeyName]string) error
	UpdateOthers(uid AccountUid, others map[KeyName]string) error
	UpdateAccountStatus(uid AccountUid, status *AccountStatus) error
}
