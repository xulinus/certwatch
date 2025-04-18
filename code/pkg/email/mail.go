package email

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"

	"github.com/xulinus/certwatch/pkg/models"
)

var (
	smtphost = os.Getenv("SMTP_HOST")
	smtpport = os.Getenv("SMTP_PORT")
	smtpuser = os.Getenv("SMTP_USER")
	smtppass = os.Getenv("SMTP_PASS")
	to       = os.Getenv("TO")
)

func SendReminder(cert models.Cert) {
	subject := fmt.Sprintf("%s expiry notification", cert.CommonName)
	body := fmt.Sprintf(`Hello! 
The certificate for %s will expire in less than 20 days (on %s). It is time to renew it now.

Kind regards,
certwatch
		`, cert.CommonName, cert.NotAfter)

	m := gomail.NewMessage()
	m.SetHeader("From", smtpuser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	port, err := strconv.Atoi(smtpport)
	if err != nil {
		log.Println("Incorrect SMTP port")
	}

	d := gomail.NewDialer(smtphost, port, smtpuser, smtppass)
	err = d.DialAndSend(m)
	if err != nil {
		log.Fatal("Error sending notification email", err)
	} else {
		log.Printf("Domain renewal remaind for %s sent to %s", cert.CommonName, to)
	}
}
