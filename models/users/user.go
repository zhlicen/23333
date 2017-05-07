package users

import (
	"23333/utils/web/beenh/beeaccount"
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

func (u *UserModel) Add(accountInfo *beeaccount.AccountInfo) error {
	o := orm.NewOrm()
	user := new(User)
	user.Uid = accountInfo.Uid
	user.Username = accountInfo.LoginIds[beeaccount.UserName.Name].Id
	user.Email = accountInfo.LoginIds[beeaccount.Email.Name].Id
	user.EmailVerified = accountInfo.LoginIds[beeaccount.Email.Name].Verified
	user.Mobile = accountInfo.LoginIds[beeaccount.Mobile.Name].Id
	user.Regtime = time.Now().Format("2006-01-02 15:04:05")
	user.Password, _ = accountInfo.Password.GetPwd()

	_, err := o.Insert(user)
	return err
}

func (u *UserModel) Delete(uid string) error {
	return errors.New("not supported")
}

func (u *UserModel) GetUserId(idName beeaccount.IdName, loginId string) (*beeaccount.UserId, error) {
	user := User{}
	var column string
	switch idName {
	case beeaccount.UserName.Name:
		column = "UserName"
		user.Username = loginId
	case beeaccount.Email.Name:
		column = "Email"
		user.Email = loginId
	case beeaccount.Mobile.Name:
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
	return &beeaccount.UserId{"", "", user.Uid}, nil
}

func (u *UserModel) GetAccountInfo(uid string) (*beeaccount.AccountInfo, error) {
	user := User{}
	user.Uid = uid

	o := orm.NewOrm()
	readErr := o.Read(&user)
	if readErr != nil {
		fmt.Println(readErr.Error())
		return nil, readErr
	}
	accountInfo := beeaccount.NewAccountInfo()
	accountInfo.Uid = user.Uid
	accountInfo.LoginIds[beeaccount.UserName.Name] = beeaccount.NewLoginId(user.Username)
	accountInfo.LoginIds[beeaccount.Email.Name] = beeaccount.NewLoginId(user.Email, user.EmailVerified)
	accountInfo.LoginIds[beeaccount.Mobile.Name] = beeaccount.NewLoginId(user.Mobile)
	accountInfo.Password.SetEncryptedPwd(user.Password)
	return accountInfo, nil
}

func (u *UserModel) GetAccountBaseInfo(uid string) (*beeaccount.AccountBaseInfo, error) {
	user := User{}
	user.Uid = uid

	o := orm.NewOrm()
	readErr := o.Read(&user)
	if readErr != nil {
		fmt.Println(readErr.Error())
		return nil, readErr
	}
	accountBaseInfo := beeaccount.NewAccountBaseInfo()
	accountBaseInfo.Uid = user.Uid
	accountBaseInfo.LoginIds[beeaccount.UserName.Name] = beeaccount.NewLoginId(user.Username)
	accountBaseInfo.LoginIds[beeaccount.Email.Name] = beeaccount.NewLoginId(user.Email)
	accountBaseInfo.LoginIds[beeaccount.Email.Name] = beeaccount.NewLoginId(user.Email, user.EmailVerified)
	accountBaseInfo.LoginIds[beeaccount.Mobile.Name] = beeaccount.NewLoginId(user.Mobile)
	accountBaseInfo.Password.SetEncryptedPwd(user.Password)
	return accountBaseInfo, nil
}

func (u *UserModel) GetOAuth2Id(uid string) (map[beeaccount.KeyName]string, error) {
	return nil, errors.New("not implemented")
}

func (u *UserModel) GetProfiles(uid string) (map[beeaccount.KeyName]string, error) {
	return nil, errors.New("not implemented")
}

func (u *UserModel) GetOthers(uid string) (map[beeaccount.KeyName]string, error) {
	return nil, errors.New("not implemented")
}
func (u *UserModel) GetAccountStatus(uid string) (*beeaccount.AccountStatus, error) {
	return nil, errors.New("not implemented")
}

func (u *UserModel) UpdateAccountBaseInfo(uid string, baseInfo *beeaccount.AccountBaseInfo) error {
	o := orm.NewOrm()
	user := new(User)
	user.Uid = baseInfo.Uid
	o.Read(user)
	user.Username = baseInfo.LoginIds[beeaccount.UserName.Name].Id
	user.Email = baseInfo.LoginIds[beeaccount.Email.Name].Id
	user.EmailVerified = baseInfo.LoginIds[beeaccount.Email.Name].Verified
	user.Mobile = baseInfo.LoginIds[beeaccount.Mobile.Name].Id
	user.Password, _ = baseInfo.Password.GetPwd()
	fmt.Println(user)
	_, err := o.Update(user)
	return err
}

func (u *UserModel) UpdateOAuth2Id(uid string, ids map[beeaccount.KeyName]string) error {
	return errors.New("not implemented")
}

func (u *UserModel) UpdateProfiles(uid string, profiles map[beeaccount.KeyName]string) error {
	return errors.New("not implemented")
}

func (u *UserModel) UpdateOthers(uid string, others map[beeaccount.KeyName]string) error {
	return errors.New("not implemented")
}

func (u *UserModel) UpdateAccountStatus(uid string, status *beeaccount.AccountStatus) error {
	return errors.New("not implemented")
}
