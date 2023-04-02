package mail

import (
	"fmt"

	"cheatppt/config"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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
	from := mail.NewEmail(subject, conf.Sender)
	to := mail.NewEmail(ctx.Username, ctx.Email)
	htmlContent := fmt.Sprintf("<p>您好，您正在进行%s邮箱验证。</p>"+
		"<p>您的验证码为: <strong>%s</strong></p>"+
		"<p>验证码 %d 分钟内有效，如果不是本人操作，请忽略。</p>",
		ctx.Username, ctx.Code, ctx.ValidMin)
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
	client := sendgrid.NewSendClient(conf.ApiKey)
	_, err := client.Send(message)

	return err
}
