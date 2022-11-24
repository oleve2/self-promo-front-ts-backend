package email

import (
	"fmt"
	"log"
	"net/smtp"
	"regexp"
	"strings"
)

type Service struct {
	EmailFrom string
	EmailTo   string
	Host      string
	Port      string
	Pass      string
}

func NewService(
	EmailFrom string,
	EmailTo string,
	Host string,
	Port string,
	Pass string,
) *Service {
	return &Service{
		EmailFrom: EmailFrom,
		EmailTo:   EmailTo,
		Host:      Host,
		Port:      Port,
		Pass:      Pass,
	}
}

// https://stackoverflow.com/questions/46369598/how-to-define-senders-name-in-golang-net-smtp-sendmail
// https://stackoverflow.com/questions/46369598/how-to-define-senders-name-in-golang-net-smtp-sendmail

func (s *Service) SendMail(emailSubject string, emailBody string) error {
	from := s.EmailFrom
	to := []string{s.EmailTo}

	text := fmt.Sprintf("Subject: %s\r\n", emailSubject) + "\r\n" + fmt.Sprintf("%s\r\n", emailBody)
	message := []byte(text)

	auth := smtp.PlainAuth("", from, s.Pass, s.Host)

	err := smtp.SendMail(s.Host+":"+s.Port, auth, from, to, message)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

/*
validate phone + email
https://stackoverflow.com/questions/45324473/regex-phone-number-using-validation-v2-golang-package-not-working
https://stackoverflow.com/questions/66624011/how-to-validate-an-email-address-in-go

https://regex101.com/
*/

func RemoveSpaces(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

func (s *Service) ValidateEmail(str string) bool {
	sTrimmed := str //RemoveSpaces(str)
	// check email
	emailRE := regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,4}$`)
	emailCheck := emailRE.MatchString(sTrimmed)
	//fmt.Println("emailCheck=", emailCheck)
	return emailCheck
}

func (s *Service) ValidatePhone(str string) bool {
	sTrimmed := RemoveSpaces(str)
	// check phone
	phoneRE := regexp.MustCompile(`^[\+0-9_\-()]{7,20}$`)
	phoneCheck := phoneRE.MatchString(sTrimmed)
	//fmt.Println("phoneCheck=", phoneCheck)
	return phoneCheck
}

func (s *Service) ValidateName(str string) bool {
	sTrimmed := RemoveSpaces(str)
	// check name
	nameRE := regexp.MustCompile(`^[A-ZА-ЯЁa-zа-яё\s]{2,30}$`)
	nameCheck := nameRE.MatchString(sTrimmed)
	//fmt.Println("nameCheck=", nameCheck)
	return nameCheck
}
