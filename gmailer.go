package gmailer

import (
	"bytes"

	gomail "gopkg.in/gomail.v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
)

type Email struct {
	Subject     string
	Body        string
	From        string
	Dest        []string
	CC          []string
	BCC         []string
	AttachFiles []string
}

type Mailer struct {
	svc sesiface.SESAPI
}

func New(svc sesiface.SESAPI) *Mailer {
	return &Mailer{svc}
}

func (m Mailer) Send(mail Email) error {
	input := m.createInput(mail)

	_, err := m.svc.SendEmail(input)
	if err != nil {
		return err

	}
	return nil
}

func (m Mailer) SendRaw(mail Email) error {
	input, err := m.createRawInput(mail)
	if err != nil {
		return err
	}

	_, err = m.svc.SendRawEmail(input)
	if err != nil {
		return err

	}
	return nil
}

func (m Mailer) createInput(mail Email) *ses.SendEmailInput {
	toAdd := convertStringSliceToAWSString(mail.Dest)
	ccAdd := convertStringSliceToAWSString(mail.CC)
	bccAdd := convertStringSliceToAWSString(mail.BCC)

	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses:  toAdd,
			CcAddresses:  ccAdd,
			BccAddresses: bccAdd,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(mail.Body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(mail.Subject),
			},
		},
		Source: aws.String(mail.From),
	}
}

func (m Mailer) createRawInput(mail Email) (*ses.SendRawEmailInput, error) {
	gm := gomail.NewMessage()
	gm.SetHeader("From", mail.From)
	gm.SetHeader("To", mail.Dest...)
	if len(mail.CC) > 0 {
		gm.SetHeader("Cc", mail.CC...)
	}
	if len(mail.CC) > 0 {
		gm.SetHeader("Bcc", mail.BCC...)
	}
	gm.SetHeader("Subject", mail.Subject)
	gm.SetBody("text/plain", mail.Body)
	for _, attachment := range mail.AttachFiles {
		gm.Attach(attachment)
	}

	var rawData bytes.Buffer
	if _, err := gm.WriteTo(&rawData); err != nil {
		return nil, err
	}

	return &ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{
			Data: rawData.Bytes(),
		},
		Source: aws.String(mail.From),
	}, nil
}

func convertStringSliceToAWSString(strs []string) (awsStr []*string) {
	for _, str := range strs {
		awsStr = append(awsStr, aws.String(str))
	}
	return
}
