# gmailer

[![Build Status](https://travis-ci.org/julianvilas/gmailer.svg?branch=master)](https://travis-ci.org/julianvilas/gmailer)

Simple Go library to send emails using [AWS SES](https://aws.amazon.com/ses/). In order to send the email you need to setup your AWS credentials as specified in the aws-sdk-go [documentation](https://github.com/aws/aws-sdk-go#configuring-credentials).

It has a function to send simple emails `SendEmail` and a function to send raw emails `SendRaw`. The latest is useful for sending emails with attachments.

# Examples

Assuming you have credentials configured in your default shared credentials profile (`~/.aws/credentials`):

```go
func main() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	}))
	svc := ses.New(sess)

	m := gmailer.New(svc)
	err := m.SendRaw(gmailer.Email{
		Subject:     "A fancy email sent with AWS SES",
		Body:        "Here I'll tell you lovely things.",
		From:        "alice@example.com",
		Dest:        "bob@example.com",
	})
	if err != nil {
		log.Panicf("executing sending email:", err)
	}
}
```
