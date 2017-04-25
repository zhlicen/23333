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
	Uid      string `orm:"pk"`
	Username string `orm:"unique"`
	Mobile   string `orm:"unique"`
	Email    string `orm:"unique"`
	Regtime  string
	Password string
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
	user.Uid = accountInfo.Uid
	user.Username = accountInfo.Ids[account.UserName.Name].Id
	user.Email = accountInfo.Ids[account.Email.Name].Id
	user.Mobile = accountInfo.Ids[account.Mobile.Name].Id
	user.Regtime = time.Now().Format("2006-01-02 15:04:05")
	user.Password, _ = accountInfo.Password.GetPwd()
	_, err := o.Insert(user)
	return err
}

func (u *UserModel) FindByPid(id string) (*account.AccountInfo, error) {
	return nil, nil
}

func (u *UserModel) FindById(id string) (*account.AccountInfo, error) {
	desc, matchErr := account.MatchIdDescriptor(id)
	if matchErr != nil {
		return nil, matchErr
	}

	user := User{}
	var column string
	switch desc.Name {
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
		return nil, errors.New("not found")
	}
	o := orm.NewOrm()
	readErr := o.Read(&user, column)
	if readErr != nil {
		fmt.Println(readErr.Error())
		return nil, readErr
	}
	accountInfo := account.NewAccountInfo()
	accountInfo.Uid = user.Uid
	accountInfo.Ids[account.UserName.Name] = account.NewAccountId(user.Username)
	accountInfo.Ids[account.Email.Name] = account.NewAccountId(user.Email)
	accountInfo.Ids[account.Mobile.Name] = account.NewAccountId(user.Mobile)
	accountInfo.Password.SetEncryptedPwd(user.Password)
	return accountInfo, nil
}
