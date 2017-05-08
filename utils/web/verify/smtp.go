package verify

import (
	"23333/utils/idgen"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"net/mail"
	"net/smtp"
	"sync"
	"time"
)

type smtpToken struct {
	token   string
	genTime time.Time
}

func (t *smtpToken) validate(timeoutSeconds uint) bool {
	timeout := t.genTime.Add(time.Duration(timeoutSeconds) * time.Second)
	if time.Now().After(timeout) {
		return false
	}
	return true
}

// SMTPVerifyService
type SMTPVerifyService struct {
	mailTitle           string
	mailTpl             string
	mailTplData         map[string]string
	tokenTimeoutSeconds uint
	tokenGenerator      idgen.IdGenerator
	tokenMap            map[string]smtpToken
	lock                sync.RWMutex
	smtpConfig          SMTPConfig
}

// NewSMTPVerifyService
func NewSMTPVerifyService(mailTpl string, tokenTimeoutSeconds uint, tokenGen idgen.IdGenerator,
	smtpConfig SMTPConfig) *SMTPVerifyService {
	vs := &SMTPVerifyService{mailTpl: mailTpl, tokenTimeoutSeconds: tokenTimeoutSeconds,
		tokenGenerator: tokenGen, smtpConfig: smtpConfig}
	vs.tokenMap = make(map[string]smtpToken)
	vs.mailTplData = make(map[string]string)
	go vs.gc()
	return vs
}

func (s *SMTPVerifyService) gc() {
	s.lock.Lock()
	defer s.lock.Unlock()
	for key, token := range s.tokenMap {
		if !token.validate(s.tokenTimeoutSeconds) {
			delete(s.tokenMap, key)
			fmt.Println("Token " + token.token + "for " + key + " expired!")
		}
	}
	time.AfterFunc(time.Duration(2)*time.Second, func() { go s.gc() })
}

func (s *SMTPVerifyService) AddMailTplData(key string, val string) {
	s.mailTplData[key] = val
}

func (s *SMTPVerifyService) SetMailTitle(title string) {
	s.mailTitle = title
}

func mergeMaps(mapTo map[string]string, maps ...interface{}) {
	for _, mapOne := range maps {
		if mapFrom, ok := mapOne.(map[string]string); ok {
			for k, v := range mapFrom {
				mapTo[k] = v
			}
		}
	}
}

func (s *SMTPVerifyService) SendToken(key string, params ...interface{}) error {
	token, err := s.tokenGenerator.Generate()
	if err != nil {
		return err
	}
	var (
		to    mail.Address
		title string
		ok    bool
	)
	if len(params) == 0 {
		return errors.New("no enough params")
	}
	if to, ok = params[0].(mail.Address); !ok {
		return errors.New("invalid param: to")
	}
	mailTitle := s.mailTitle
	if title, ok = params[1].(string); ok {
		if title != "" {
			mailTitle = title
		}
	}

	data := make(map[string]string)
	mergeMaps(data, s.mailTplData)
	if len(params) > 2 {
		if dataMap, ok := params[2].(map[string]string); ok {
			mergeMaps(data, dataMap)
		}
	}

	data["token"] = token
	content, mailErr := s.genMailContent(data)
	if mailErr != nil {
		return mailErr
	}
	go s.sendMail(mailTitle, to, content)

	{
		s.lock.Lock()
		defer s.lock.Unlock()
		s.tokenMap[key] = smtpToken{token, time.Now()}
		fmt.Println("Token " + token + "for " + key + " added!")
	}
	return nil

}

func (s *SMTPVerifyService) Verify(key string, token string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if savedToken, ok := s.tokenMap[key]; ok {
		if savedToken.token == token {
			delete(s.tokenMap, key)
			fmt.Println("Token " + token + "for " + key + " verified and deleted!")
			return nil
		}
	}
	return errors.New("invalid token")
}

func (s *SMTPVerifyService) sendMail(title string, to mail.Address, content string) {
	config := s.smtpConfig

	auth := smtp.PlainAuth(
		"",
		config.Account,
		config.Password,
		config.SMTPServer,
	)
	header := make(map[string]string)
	header["From"] = config.From.String()
	header["To"] = to.String()
	header["Subject"] = title
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(content))
	mailErr := smtp.SendMail(config.SMTPServer+":25", auth, config.From.Address,
		[]string{to.Address}, []byte(message))
	if mailErr != nil {
		fmt.Println("Mail send failed:" + mailErr.Error())
		return
	}
	fmt.Println("Mail sent to:" + to.Address)
}

func (s *SMTPVerifyService) genMailContent(data interface{}) (string, error) {
	t, err := template.ParseFiles(s.mailTpl)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

type SMTPConfig struct {
	SMTPServer string
	Account    string
	Password   string
	From       mail.Address
}
