package gmailer

import (
	"bytes"

	gomail "gopkg.in/gomail.v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
)

// Email defines the structure of an email to be sent.
type Email struct {
	Subject     string
	Body        string
	From        string
	Dest        []string
	CC          []string
	BCC         []string
	AttachFiles []string
	ConfigSet   string // AWS SES Configuration Set (https://goo.gl/sXr7fj).
	HTML        bool   // Specifies whether the content of the email is HTML or Text.
}

// Mailer is a structure that contains a service that implements the AWS SES interface.
type Mailer struct {
	svc sesiface.SESAPI
}

// New receives a service that implements the AWS ses interface and returns a Mailer associated to it.
func New(svc sesiface.SESAPI) *Mailer {
	return &Mailer{svc}
}

// Send sends the Email.
func (m Mailer) Send(mail Email) error {
	input := m.createInput(mail)

	_, err := m.svc.SendEmail(input)
	if err != nil {
		return err

	}
	return nil
}

// SendRaw sends the email using RawInput, useful when sending attached files.
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

	content := &ses.Content{
		Charset: aws.String("UTF-8"),
		Data:    aws.String(mail.Body),
	}
	var body *ses.Body
	if mail.HTML {
		body = &ses.Body{Html: content}
	} else {
		body = &ses.Body{Text: content}
	}

	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses:  toAdd,
			CcAddresses:  ccAdd,
			BccAddresses: bccAdd,
		},
		Message: &ses.Message{
			Body: body,
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(mail.Subject),
			},
		},
		Source:               aws.String(mail.From),
		ConfigurationSetName: aws.String(mail.ConfigSet),
	}
}

func (m Mailer) createRawInput(mail Email) (*ses.SendRawEmailInput, error) {
	gm := gomail.NewMessage()
	gm.SetHeader("From", mail.From)
	gm.SetHeader("To", mail.Dest...)
	if len(mail.CC) > 0 {
		gm.SetHeader("Cc", mail.CC...)
	}
	if len(mail.BCC) > 0 {
		gm.SetHeader("Bcc", mail.BCC...)
	}
	gm.SetHeader("Subject", mail.Subject)

	if mail.ConfigSet != "" {
		gm.SetHeader("X-SES-CONFIGURATION-SET", mail.ConfigSet)
	}

	var contentType string
	if mail.HTML {
		contentType = "text/html;charset=UTF-8"
	} else {
		contentType = "text/plain;charset=UTF-8"
	}
	gm.SetBody(contentType, mail.Body)

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
