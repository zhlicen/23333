package users

import (
	"23333/utils/web/beenhance/beaccount"
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

func (u *UserModel) Add(accountInfo *beaccount.AccountInfo) error {
	o := orm.NewOrm()
	user := new(User)
	user.Uid = accountInfo.Uid.String()
	user.Username = accountInfo.Ids[beaccount.UserName.Name].Id
	user.Email = accountInfo.Ids[beaccount.Email.Name].Id
	user.EmailVerified = accountInfo.Ids[beaccount.Email.Name].Verified
	user.Mobile = accountInfo.Ids[beaccount.Mobile.Name].Id
	user.Regtime = time.Now().Format("2006-01-02 15:04:05")
	user.Password, _ = accountInfo.Password.GetPwd()

	_, err := o.Insert(user)
	return err
}

func (u *UserModel) Delete(uid beaccount.AccountUid) error {
	return errors.New("not supported")
}

func (u *UserModel) GetUidById(idName beaccount.IdName, id string) (beaccount.AccountUid, error) {
	user := User{}
	var column string
	switch idName {
	case beaccount.UserName.Name:
		column = "UserName"
		user.Username = id
	case beaccount.Email.Name:
		column = "Email"
		user.Email = id
	case beaccount.Mobile.Name:
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
	return beaccount.AccountUid(user.Uid), nil
}

func (u *UserModel) GetAccountInfo(uid beaccount.AccountUid) (*beaccount.AccountInfo, error) {
	user := User{}
	user.Uid = uid.String()

	o := orm.NewOrm()
	readErr := o.Read(&user)
	if readErr != nil {
		fmt.Println(readErr.Error())
		return nil, readErr
	}
	accountInfo := beaccount.NewAccountInfo()
	accountInfo.Uid.SetVal(user.Uid)
	accountInfo.Ids[beaccount.UserName.Name] = beaccount.NewAccountId(user.Username)
	accountInfo.Ids[beaccount.Email.Name] = beaccount.NewAccountId(user.Email, user.EmailVerified)
	accountInfo.Ids[beaccount.Mobile.Name] = beaccount.NewAccountId(user.Mobile)
	accountInfo.Password.SetEncryptedPwd(user.Password)
	return accountInfo, nil
}

func (u *UserModel) GetAccountBaseInfo(uid beaccount.AccountUid) (*beaccount.AccountBaseInfo, error) {
	user := User{}
	user.Uid = uid.String()

	o := orm.NewOrm()
	readErr := o.Read(&user)
	if readErr != nil {
		fmt.Println(readErr.Error())
		return nil, readErr
	}
	accountBaseInfo := beaccount.NewAccountBaseInfo()
	accountBaseInfo.Uid.SetVal(user.Uid)
	accountBaseInfo.Ids[beaccount.UserName.Name] = beaccount.NewAccountId(user.Username)
	accountBaseInfo.Ids[beaccount.Email.Name] = beaccount.NewAccountId(user.Email)
	accountBaseInfo.Ids[beaccount.Email.Name] = beaccount.NewAccountId(user.Email, user.EmailVerified)
	accountBaseInfo.Ids[beaccount.Mobile.Name] = beaccount.NewAccountId(user.Mobile)
	accountBaseInfo.Password.SetEncryptedPwd(user.Password)
	return accountBaseInfo, nil
}

func (u *UserModel) GetOAuth2Id(uid beaccount.AccountUid) (map[beaccount.KeyName]string, error) {
	return nil, errors.New("not implemented")
}

func (u *UserModel) GetProfiles(uid beaccount.AccountUid) (map[beaccount.KeyName]string, error) {
	return nil, errors.New("not implemented")
}

func (u *UserModel) GetOthers(uid beaccount.AccountUid) (map[beaccount.KeyName]string, error) {
	return nil, errors.New("not implemented")
}
func (u *UserModel) GetAccountStatus(uid beaccount.AccountUid) (*beaccount.AccountStatus, error) {
	return nil, errors.New("not implemented")
}

func (u *UserModel) UpdateAccountBaseInfo(uid beaccount.AccountUid, baseInfo *beaccount.AccountBaseInfo) error {
	o := orm.NewOrm()
	user := new(User)
	user.Uid = baseInfo.Uid.String()
	o.Read(user)
	user.Username = baseInfo.Ids[beaccount.UserName.Name].Id
	user.Email = baseInfo.Ids[beaccount.Email.Name].Id
	user.EmailVerified = baseInfo.Ids[beaccount.Email.Name].Verified
	user.Mobile = baseInfo.Ids[beaccount.Mobile.Name].Id
	user.Password, _ = baseInfo.Password.GetPwd()
	fmt.Println(user)
	_, err := o.Update(user)
	return err
}

func (u *UserModel) UpdateOAuth2Id(uid beaccount.AccountUid, ids map[beaccount.KeyName]string) error {
	return errors.New("not implemented")
}

func (u *UserModel) UpdateProfiles(uid beaccount.AccountUid, profiles map[beaccount.KeyName]string) error {
	return errors.New("not implemented")
}

func (u *UserModel) UpdateOthers(uid beaccount.AccountUid, others map[beaccount.KeyName]string) error {
	return errors.New("not implemented")
}

func (u *UserModel) UpdateAccountStatus(uid beaccount.AccountUid, status *beaccount.AccountStatus) error {
	return errors.New("not implemented")
}
