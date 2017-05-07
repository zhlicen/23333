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
	user.Uid = accountInfo.Uid
	user.Username = accountInfo.LoginIds[beaccount.UserName.Name].Id
	user.Email = accountInfo.LoginIds[beaccount.Email.Name].Id
	user.EmailVerified = accountInfo.LoginIds[beaccount.Email.Name].Verified
	user.Mobile = accountInfo.LoginIds[beaccount.Mobile.Name].Id
	user.Regtime = time.Now().Format("2006-01-02 15:04:05")
	user.Password, _ = accountInfo.Password.GetPwd()

	_, err := o.Insert(user)
	return err
}

func (u *UserModel) Delete(uid string) error {
	return errors.New("not supported")
}

func (u *UserModel) GetUserId(idName beaccount.IdName, loginId string) (*beaccount.UserId, error) {
	user := User{}
	var column string
	switch idName {
	case beaccount.UserName.Name:
		column = "UserName"
		user.Username = loginId
	case beaccount.Email.Name:
		column = "Email"
		user.Email = loginId
	case beaccount.Mobile.Name:
		column = "Mobile"
		user.Mobile = loginId
	default:
		return nil, errors.New("not found")
	}
	fmt.Println(column + "::" + loginId)
	o := orm.NewOrm()
	readErr := o.Read(&user, column)
	if readErr != nil {
		return nil, readErr
	}
	return &beaccount.UserId{"", "", user.Uid}, nil
}

func (u *UserModel) GetAccountInfo(uid string) (*beaccount.AccountInfo, error) {
	user := User{}
	user.Uid = uid

	o := orm.NewOrm()
	readErr := o.Read(&user)
	if readErr != nil {
		fmt.Println(readErr.Error())
		return nil, readErr
	}
	accountInfo := beaccount.NewAccountInfo()
	accountInfo.Uid = user.Uid
	accountInfo.LoginIds[beaccount.UserName.Name] = beaccount.NewLoginId(user.Username)
	accountInfo.LoginIds[beaccount.Email.Name] = beaccount.NewLoginId(user.Email, user.EmailVerified)
	accountInfo.LoginIds[beaccount.Mobile.Name] = beaccount.NewLoginId(user.Mobile)
	accountInfo.Password.SetEncryptedPwd(user.Password)
	return accountInfo, nil
}

func (u *UserModel) GetAccountBaseInfo(uid string) (*beaccount.AccountBaseInfo, error) {
	user := User{}
	user.Uid = uid

	o := orm.NewOrm()
	readErr := o.Read(&user)
	if readErr != nil {
		fmt.Println(readErr.Error())
		return nil, readErr
	}
	accountBaseInfo := beaccount.NewAccountBaseInfo()
	accountBaseInfo.Uid = user.Uid
	accountBaseInfo.LoginIds[beaccount.UserName.Name] = beaccount.NewLoginId(user.Username)
	accountBaseInfo.LoginIds[beaccount.Email.Name] = beaccount.NewLoginId(user.Email)
	accountBaseInfo.LoginIds[beaccount.Email.Name] = beaccount.NewLoginId(user.Email, user.EmailVerified)
	accountBaseInfo.LoginIds[beaccount.Mobile.Name] = beaccount.NewLoginId(user.Mobile)
	accountBaseInfo.Password.SetEncryptedPwd(user.Password)
	return accountBaseInfo, nil
}

func (u *UserModel) GetOAuth2Id(uid string) (map[beaccount.KeyName]string, error) {
	return nil, errors.New("not implemented")
}

func (u *UserModel) GetProfiles(uid string) (map[beaccount.KeyName]string, error) {
	return nil, errors.New("not implemented")
}

func (u *UserModel) GetOthers(uid string) (map[beaccount.KeyName]string, error) {
	return nil, errors.New("not implemented")
}
func (u *UserModel) GetAccountStatus(uid string) (*beaccount.AccountStatus, error) {
	return nil, errors.New("not implemented")
}

func (u *UserModel) UpdateAccountBaseInfo(uid string, baseInfo *beaccount.AccountBaseInfo) error {
	o := orm.NewOrm()
	user := new(User)
	user.Uid = baseInfo.Uid
	o.Read(user)
	user.Username = baseInfo.LoginIds[beaccount.UserName.Name].Id
	user.Email = baseInfo.LoginIds[beaccount.Email.Name].Id
	user.EmailVerified = baseInfo.LoginIds[beaccount.Email.Name].Verified
	user.Mobile = baseInfo.LoginIds[beaccount.Mobile.Name].Id
	user.Password, _ = baseInfo.Password.GetPwd()
	fmt.Println(user)
	_, err := o.Update(user)
	return err
}

func (u *UserModel) UpdateOAuth2Id(uid string, ids map[beaccount.KeyName]string) error {
	return errors.New("not implemented")
}

func (u *UserModel) UpdateProfiles(uid string, profiles map[beaccount.KeyName]string) error {
	return errors.New("not implemented")
}

func (u *UserModel) UpdateOthers(uid string, others map[beaccount.KeyName]string) error {
	return errors.New("not implemented")
}

func (u *UserModel) UpdateAccountStatus(uid string, status *beaccount.AccountStatus) error {
	return errors.New("not implemented")
}
