package account

type AccountModel interface {
	Add(accountInfo *AccountInfo) error
	//Delete(uid string) error
	FindByUid(id string) (*AccountInfo, error)
	FindById(id string) (*AccountInfo, error)
}
