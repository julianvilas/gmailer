# gmailer

Simple Go library to send emails using AWS SES.

## Examples

```go
func main() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1"),
	}))
	svc := ses.New(sess)

	m := gmailer.New(svc)
	err := m.SendRaw(mailer.Email{
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
