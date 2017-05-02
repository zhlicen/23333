package users

import (
	"23333/account"
	"encoding/base64"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/astaxie/beego"
)

type RegisterController struct {
	beego.Controller
}

func (c *RegisterController) Get() {
	if c.GetSession("LoginUser") != nil {
		fmt.Println("Logout current user to register")
		c.DelSession("LoginUser")
	}
	c.TplName = "register/index.tpl"
}

func (c *RegisterController) Post() {
	accountInfo := account.NewAccountInfo()
	username := c.GetString("username")
	password := c.GetString("password")
	mobile := c.GetString("mobile")
	email := c.GetString("email")
	accountInfo.Domain = beego.BConfig.AppName
	accountInfo.Ids[account.UserName.Name] = account.NewAccountId(username)
	fmt.Println("UserName:" + username)
	accountInfo.Ids[account.Mobile.Name] = account.NewAccountId(mobile)
	fmt.Println("Mobile:" + mobile)
	accountInfo.Ids[account.Email.Name] = account.NewAccountId(email)
	fmt.Println("Email:" + email)

	accountInfo.Password.SetPwd("Password", password, accountInfo.Uid, pwdEncryptorSalt)
	pwd, err := accountInfo.Password.GetPwd()
	if err == nil {
		fmt.Println("Password:" + pwd)
	} else {
		fmt.Println(err)
	}

	regErr := accountService.Register(c.Ctx, accountInfo)
	if regErr != nil {
		fmt.Println(regErr.Error())
		c.Ctx.Redirect(302, "/register")
		return
	}
	// go SendMail(email)
	c.Ctx.Redirect(302, "/login")
}

func encodeRFC2047(String string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>")
}

func SendMail(addr string) {
	smtpServer := "smtp.sohu.com"
	from := mail.Address{"3j noreply", "jjj_noreply@sohu.com"}
	to := mail.Address{"smtp", addr}
	auth := smtp.PlainAuth(
		"",
		"jjj_noreply@sohu.com",
		"232323",
		smtpServer,
	)
	title := "23333 Registration"
	body := "Register Success!"

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = encodeRFC2047(title)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))
	mailErr := smtp.SendMail(smtpServer+":25", auth, from.Address,
		[]string{to.Address}, []byte(message))
	if mailErr != nil {
		fmt.Println("Mail send failed:" + mailErr.Error())
		return
	}
	fmt.Println("Mail sent to:" + addr)
}
