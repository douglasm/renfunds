package mail

import (
	"fmt"
	"log"

	"gopkg.in/mailgun/mailgun-go.v1"

	"ssafa/crypto"
)

const (
	domain        = "renfunds.barcodeprograms.co.uk"
	mailGunDomain = "91e1852a0b2d67a4bd8b1a5b500655102f3e7560e11a1a814c354d38c03718431b55351fe82cff905f57e63f79440321"
	mailKey       = "9dcd1829f0f74387f7863d55bc6b9a4cd36c25ad911f874df658614466e815b0295c3e2784fddea3c04bda4942ad18be850db1454074322c4164cbdcbca0a3a8"
)

// SendActivate will use mailgun to send a profile activate e-mail.
func SendActivate(email, code string) {
	mgDomain := crypto.Decrypt(mailGunDomain)
	mgKey := crypto.Decrypt(mailKey)

	messageText := fmt.Sprintf("Please go to https://%s/activate/%s", domain, code)
	messageText += " to activate your profile, or go to "
	messageText += fmt.Sprintf("https://%s/activate and enter the code %s", domain, code)

	sendAddress := fmt.Sprintf("noreply@%s", domain)

	mg := mailgun.NewMailgun(mgDomain, mgKey, "")
	m := mg.NewMessage(
		sendAddress,
		"Welcome to Renfunds",
		messageText,
		email,
	)

	_, _, err := mg.Send(m)
	if err != nil {
		log.Println("Error: send activate", err)
	}
}

// SendReset will use mailgun to send a password reset e-mail.
func SendReset(email, code string) {
	mgDomain := crypto.Decrypt(mailGunDomain)
	mgKey := crypto.Decrypt(mailKey)

	messageText := fmt.Sprintf("Please go to https://%s/reset/%s", domain, code)
	messageText += " to reset your password, or go to "
	messageText += fmt.Sprintf("https://%s/reset and enter the code %s", domain, code)

	sendAddress := fmt.Sprintf("noreply@%s", domain)

	mg := mailgun.NewMailgun(mgDomain, mgKey, "")
	m := mg.NewMessage(
		sendAddress,
		"Renfunds password reset",
		messageText,
		email,
	)

	_, _, err := mg.Send(m)
	if err != nil {
		log.Println("Error: send reset", err)
	}
}
