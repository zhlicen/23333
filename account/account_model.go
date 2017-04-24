package account

type AccountModel interface {
	Add(accountInfo *AccountInfo) error
	FindByPid(id string) (*AccountInfo, error)
	FindById(id string) (*AccountInfo, error)
}
