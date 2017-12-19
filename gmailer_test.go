package gmailer

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
)

var tcSend = []struct {
	name             string
	skip, skipAlways bool
	mock             sesiface.SESAPI
	mail             Email
	wantErr          bool
}{
	{
		name: "text",
		mock: sesMock{},
		mail: Email{
			Subject: "Test subject",
			Body:    "Test body",
			From:    "from@example.com",
			Dest:    []string{"to@example.com"},
		},
		wantErr: false,
	},
	{
		name: "html",
		mock: sesMock{},
		mail: Email{
			Subject: "Test subject",
			Body:    "Test body",
			From:    "from@example.com",
			Dest:    []string{"to@example.com"},
			HTML:    true,
		},
		wantErr: false,
	},
	{
		name: "cc and bcc",
		mock: sesMock{},
		mail: Email{
			Subject: "Test subject",
			Body:    "Test body",
			From:    "from@example.com",
			Dest:    []string{"to@example.com"},
			CC:      []string{"cc@example.com"},
			BCC:     []string{"bcc@example.com"},
		},
		wantErr: false,
	},
	{
		name: "configuration set",
		mock: sesMock{},
		mail: Email{
			Subject:   "Test subject",
			Body:      "Test body",
			From:      "from@example.com",
			Dest:      []string{"to@example.com"},
			ConfigSet: "a fake config set",
		},
		wantErr: false,
	},
	{
		name: "returns error",
		mock: sesMock{err: errors.New("a creepy error")},
		mail: Email{
			Subject:   "Test subject",
			Body:      "Test body",
			From:      "from@example.com",
			Dest:      []string{"to@example.com"},
			ConfigSet: "a fake config set",
		},
		wantErr: true,
	},
	{
		name: "attachment",
		mock: sesMock{},
		mail: Email{
			Subject:     "Test subject",
			Body:        "Test body",
			From:        "from@example.com",
			Dest:        []string{"to@example.com"},
			AttachFiles: []string{"testdata/attachment.txt"},
		},
		wantErr: false,
	},
}

func TestSend(t *testing.T) {
	// Test all the test cases defined in tcSend
	for _, tc := range tcSend {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			if testing.Short() && tc.skip || tc.skipAlways {
				t.SkipNow()
			}

			m := New(tc.mock)

			err := m.Send(tc.mail)

			if err == nil && tc.wantErr {
				t.Error("wants error, got nil")
			} else if err != nil && !tc.wantErr {
				t.Errorf("wants nill error, got %v", err)
			}
		})
	}
}

func TestSendRaw(t *testing.T) {
	// Test all the test cases defined in tcSend
	for _, tc := range tcSend {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			if testing.Short() && tc.skip || tc.skipAlways {
				t.SkipNow()
			}

			m := New(tc.mock)

			err := m.SendRaw(tc.mail)

			if err == nil && tc.wantErr {
				t.Error("wants error, got nil")
			} else if err != nil && !tc.wantErr {
				t.Errorf("wants nill error, got %v", err)
			}
		})
	}
}

type sesMock struct {
	sesiface.SESAPI
	emailOut    *ses.SendEmailOutput
	rawEmailOut *ses.SendRawEmailOutput
	err         error
}

func (m sesMock) SendRawEmail(*ses.SendRawEmailInput) (*ses.SendRawEmailOutput, error) {
	return m.rawEmailOut, m.err
}

func (m sesMock) SendEmail(*ses.SendEmailInput) (*ses.SendEmailOutput, error) {
	return m.emailOut, m.err
}
