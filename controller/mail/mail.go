package mail

import (
	"cheatppt/config"
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func EmailVerificationSend(username string, email string) error {
	conf := config.GlobalCfg.Mail
	from := mail.NewEmail("No Reply", "noreply@example.com")
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail(username, email)
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := `
	Hi [name],
	
	Thanks for getting started with our [customer portal]!
	
	We need a little more information to complete your registration, including a confirmation of your email address.
	
	Click below to confirm your email address:
	
	[link]
	
	If you have problems, please paste the above URL into your web browser.`
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(conf.ApiKey)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	return nil
}

func PasswdResetSend() {
	conf := config.GlobalCfg.Mail
	_ = sendgrid.GetRequest(conf.ApiKey, "/v3/mail/send")
	return
}
