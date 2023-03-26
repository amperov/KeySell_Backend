package mailer

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
}

func (m *Mailer) SendNewPassword(to string, password string) {

	dialer := gomail.NewDialer("smtp.yandex.ru", 25, "alex.amperov", "UnflatIsSad")
	dial, err := dialer.Dial()
	if err != nil {
		logrus.Print(err)
		return
	}
	message := gomail.NewMessage()
	message.SetHeaders(map[string][]string{
		"From":    {message.FormatAddress("", "keys-store.online")},
		"To":      {to},
		"Subject": {"Новый пароль"}})

	message.SetBody("text/html", fmt.Sprintf("<h1>Ваш новый пароль: %s</h1>", password))
	err = dial.Send("alex.amperov@yandex.ru", []string{to}, message)
	if err != nil {
		logrus.Print(err)
		return
	}
}
