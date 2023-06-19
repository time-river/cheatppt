package mail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"

	log "github.com/sirupsen/logrus"

	"cheatppt/config"
)

type CodeCtx struct {
	Username string
	Email    string
	Code     string
	ValidMin int
}

func SendCode(ctx *CodeCtx) error {
	conf := config.Mail

	subject := "邮箱验证邮件"
	htmlContent := fmt.Sprintf("<p>您好，您正在进行%s邮箱验证。</p>"+
		"<p>您的验证码为: <strong>%s</strong></p>"+
		"<p>验证码 %d 分钟内有效，如果不是本人操作，请忽略。</p>",
		ctx.Username, ctx.Code, ctx.ValidMin)

	// setup header
	from := mail.Address{Name: conf.SenderName, Address: conf.SenderAddr}
	to := mail.Address{Name: ctx.Username, Address: ctx.Email}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlContent

	mail := &TLSMail{
		serverName: conf.SMTPServer,
		senderAddr: conf.SenderAddr,
		secret:     conf.Secret,

		fromAddr: from.Address,
		toAddr:   to.Address,
		message:  message,
	}
	return tlsSendMail(mail)
}

type TLSMail struct {
	serverName string
	senderAddr string
	secret     string

	fromAddr string
	toAddr   string
	message  string
}

func tlsSendMail(mail *TLSMail) error {
	// Connect to the SMTP Server
	host, _, _ := net.SplitHostPort(mail.serverName)

	auth := smtp.PlainAuth("", mail.senderAddr, mail.secret, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", mail.serverName, tlsconfig)
	if err != nil {
		log.Warnf("SendMail ERROR: %s\n", err.Error())
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Warnf("SendMail ERROR: %s\n", err.Error())
		return err
	}
	defer c.Quit()

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Warnf("SendMail ERROR: %s\n", err.Error())
		return err
	}

	// To && From
	if err = c.Mail(mail.fromAddr); err != nil {
		log.Warnf("SendMail ERROR: %s\n", err.Error())
		return err
	}

	if err = c.Rcpt(mail.toAddr); err != nil {
		log.Warnf("SendMail ERROR: %s\n", err.Error())
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Warnf("SendMail ERROR: %s\n", err.Error())
		return err
	}
	defer w.Close()

	_, err = w.Write([]byte(mail.message))
	if err != nil {
		log.Warnf("SendMail ERROR: %s\n", err.Error())
		return err
	}

	return nil
}
