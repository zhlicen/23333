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
	user.Username = accountInfo.LoginIDs["UserName"].ID
	user.Email = accountInfo.LoginIDs["Email"].ID
	user.EmailVerified = accountInfo.LoginIDs["Email"].Verified
	user.Mobile = accountInfo.LoginIDs["Mobile"].ID
	user.Regtime = time.Now().Format("2006-01-02 15:04:05")
	user.Password, _ = accountInfo.Password.GetPwd()

	_, err := o.Insert(user)
	return err
}

func (u *UserModel) Delete(uid string) error {
	return errors.New("not supported")
}

func (u *UserModel) GetUserID(idName string, loginID string) (*beeaccount.UserID, error) {
	user := User{}
	var column string
	switch idName {
	case "UserName":
		column = "UserName"
		user.Username = loginID
	case "Email":
		column = "Email"
		user.Email = loginID
	case "Mobile":
		column = "Mobile"
		user.Mobile = loginID
	default:
		return nil, errors.New(idName + "not found")
	}
	fmt.Println(column + "::" + loginID)
	o := orm.NewOrm()
	readErr := o.Read(&user, column)
	if readErr != nil {
		return nil, readErr
	}
	return &beeaccount.UserID{"", "", user.Uid}, nil
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
	accountInfo, _ := beeaccount.NewAccountInfo("23333")
	accountInfo.Uid = user.Uid
	accountInfo.LoginIDs["UserName"] = beeaccount.NewLoginID(user.Username)
	accountInfo.LoginIDs["Email"] = beeaccount.NewLoginID(user.Email, user.EmailVerified)
	accountInfo.LoginIDs["Mobile"] = beeaccount.NewLoginID(user.Mobile)
	accountInfo.Password.SetEncryptedPwd(user.Password)
	return accountInfo, nil
}

func (u *UserModel) GetAccountBasicInfo(uid string) (*beeaccount.AccountBasicInfo, error) {
	user := User{}
	user.Uid = uid

	o := orm.NewOrm()
	readErr := o.Read(&user)
	if readErr != nil {
		fmt.Println(readErr.Error())
		return nil, readErr
	}
	accountBaseInfo := beeaccount.NewAccountBasicInfo("23333")
	accountBaseInfo.Uid = user.Uid
	accountBaseInfo.LoginIDs["UserName"] = beeaccount.NewLoginID(user.Username)
	accountBaseInfo.LoginIDs["Email"] = beeaccount.NewLoginID(user.Email)
	accountBaseInfo.LoginIDs["Email"] = beeaccount.NewLoginID(user.Email, user.EmailVerified)
	accountBaseInfo.LoginIDs["Mobile"] = beeaccount.NewLoginID(user.Mobile)
	accountBaseInfo.Password.SetEncryptedPwd(user.Password)
	return accountBaseInfo, nil
}

func (u *UserModel) GetOAuth2ID(uid string) (map[string]string, error) {
	return nil, errors.New("not implemented")
}

func (u *UserModel) GetProfiles(uid string) (map[string]string, error) {
	return nil, errors.New("not implemented")
}

func (u *UserModel) GetOthers(uid string) (map[string]string, error) {
	return nil, errors.New("not implemented")
}
func (u *UserModel) GetAccountStatus(uid string) (*beeaccount.AccountStatus, error) {
	return nil, errors.New("not implemented")
}

func (u *UserModel) UpdateAccountBasicInfo(uid string, baseInfo *beeaccount.AccountBasicInfo) error {
	o := orm.NewOrm()
	user := new(User)
	user.Uid = baseInfo.Uid
	o.Read(user)
	user.Username = baseInfo.LoginIDs["UserName"].ID
	user.Email = baseInfo.LoginIDs["Email"].ID
	user.EmailVerified = baseInfo.LoginIDs["Email"].Verified
	user.Mobile = baseInfo.LoginIDs["Mobile"].ID
	user.Password, _ = baseInfo.Password.GetPwd()
	fmt.Println(user)
	_, err := o.Update(user)
	return err
}

func (u *UserModel) UpdateOAuth2ID(uid string, ids map[string]string) error {
	return errors.New("not implemented")
}

func (u *UserModel) UpdateProfiles(uid string, profiles map[string]string) error {
	return errors.New("not implemented")
}

func (u *UserModel) UpdateOthers(uid string, others map[string]string) error {
	return errors.New("not implemented")
}

func (u *UserModel) UpdateAccountStatus(uid string, status *beeaccount.AccountStatus) error {
	return errors.New("not implemented")
}
