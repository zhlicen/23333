package users

import (
	"23333/account"
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Uid           string `orm:"pk"`
	Username      string `orm:"unique"`
	Mobile        string `orm:"unique"`
	Email         string `orm:"unique"`
	EmailVerified bool
	Regtime       string
	Password      string
}

func (u *User) TableName() string {
	return "23333_user"
}

func init() {
	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("mysql registered")
	}
	orm.RegisterDataBase("default", "mysql", "root:125801@tcp/app_23333?charset=utf8", 30)
	orm.RegisterModel(new(User))
}

type UserModel struct {
}

func (u *UserModel) Add(accountInfo *account.AccountInfo) error {
	o := orm.NewOrm()
	user := new(User)
	user.Uid = accountInfo.Uid.String()
	user.Username = accountInfo.Ids[account.UserName.Name].Id
	user.Email = accountInfo.Ids[account.Email.Name].Id
	user.EmailVerified = accountInfo.Ids[account.Email.Name].Verified
	user.Mobile = accountInfo.Ids[account.Mobile.Name].Id
	user.Regtime = time.Now().Format("2006-01-02 15:04:05")
	user.Password, _ = accountInfo.Password.GetPwd()

	_, err := o.Insert(user)
	return err
}

func (u *UserModel) Delete(uid account.AccountUid) error {
	return errors.New("not supported")
}

func (u *UserModel) GetUidById(idName account.IdName, id string) (account.AccountUid, error) {
	user := User{}
	var column string
	switch idName {
	case account.UserName.Name:
		column = "UserName"
		user.Username = id
	case account.Email.Name:
		column = "Email"
		user.Email = id
	case account.Mobile.Name:
		column = "Mobile"
		user.Mobile = id
	default:
		return "", errors.New("not found")
	}
	fmt.Println(column + "::" + id)
	o := orm.NewOrm()
	readErr := o.Read(&user, column)
	if readErr != nil {
		return "", readErr
	}
	return account.AccountUid(user.Uid), nil
}

func (u *UserModel) GetAccountInfo(uid account.AccountUid) (*account.AccountInfo, error) {
	user := User{}
	user.Uid = uid.String()

	o := orm.NewOrm()
	readErr := o.Read(&user)
	if readErr != nil {
		fmt.Println(readErr.Error())
		return nil, readErr
	}
	accountInfo := account.NewAccountInfo()
	accountInfo.Uid.SetVal(user.Uid)
	accountInfo.Ids[account.UserName.Name] = account.NewAccountId(user.Username)
	accountInfo.Ids[account.Email.Name] = account.NewAccountId(user.Email, user.EmailVerified)
	accountInfo.Ids[account.Mobile.Name] = account.NewAccountId(user.Mobile)
	accountInfo.Password.SetEncryptedPwd(user.Password)
	return accountInfo, nil
}

func (u *UserModel) GetAccountBaseInfo(uid account.AccountUid) (*account.AccountBaseInfo, error) {
	user := User{}
	user.Uid = uid.String()

	o := orm.NewOrm()
	readErr := o.Read(&user)
	if readErr != nil {
		fmt.Println(readErr.Error())
		return nil, readErr
	}
	accountBaseInfo := account.NewAccountBaseInfo()
	accountBaseInfo.Uid.SetVal(user.Uid)
	accountBaseInfo.Ids[account.UserName.Name] = account.NewAccountId(user.Username)
	accountBaseInfo.Ids[account.Email.Name] = account.NewAccountId(user.Email)
	accountBaseInfo.Ids[account.Email.Name] = account.NewAccountId(user.Email, user.EmailVerified)
	accountBaseInfo.Ids[account.Mobile.Name] = account.NewAccountId(user.Mobile)
	accountBaseInfo.Password.SetEncryptedPwd(user.Password)
	return accountBaseInfo, nil
}

func (u *UserModel) GetOAuth2Id(uid account.AccountUid) (map[account.KeyName]string, error) {
	return nil, errors.New("not implemented")
}

func (u *UserModel) GetProfiles(uid account.AccountUid) (map[account.KeyName]string, error) {
	return nil, errors.New("not implemented")
}

func (u *UserModel) GetOthers(uid account.AccountUid) (map[account.KeyName]string, error) {
	return nil, errors.New("not implemented")
}
func (u *UserModel) GetAccountStatus(uid account.AccountUid) (*account.AccountStatus, error) {
	return nil, errors.New("not implemented")
}

func (u *UserModel) UpdateAccountBaseInfo(uid account.AccountUid, baseInfo *account.AccountBaseInfo) error {
	o := orm.NewOrm()
	user := new(User)
	user.Uid = baseInfo.Uid.String()
	o.Read(user)
	user.Username = baseInfo.Ids[account.UserName.Name].Id
	user.Email = baseInfo.Ids[account.Email.Name].Id
	user.EmailVerified = baseInfo.Ids[account.Email.Name].Verified
	user.Mobile = baseInfo.Ids[account.Mobile.Name].Id
	user.Password, _ = baseInfo.Password.GetPwd()
	fmt.Println(user)
	_, err := o.Update(user)
	return err
}

func (u *UserModel) UpdateOAuth2Id(uid account.AccountUid, ids map[account.KeyName]string) error {
	return errors.New("not implemented")
}

func (u *UserModel) UpdateProfiles(uid account.AccountUid, profiles map[account.KeyName]string) error {
	return errors.New("not implemented")
}

func (u *UserModel) UpdateOthers(uid account.AccountUid, others map[account.KeyName]string) error {
	return errors.New("not implemented")
}

func (u *UserModel) UpdateAccountStatus(uid account.AccountUid, status *account.AccountStatus) error {
	return errors.New("not implemented")
}
