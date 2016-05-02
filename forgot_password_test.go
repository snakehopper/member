package member

import (
	"fmt"
	"testing"
)

const (
	tUserEmail = "user@example.com"
)

func TestForgotPasswordCtrl_ForgotPassword(t *testing.T) {
	m := newMemberFinderTD()
	tk := newTokenGeneratorTD()
	e := newEmailerTD()
	ctrl := NewForgotPasswordCtrl(m, tk, &e)

	if err := ctrl.ForgotPassword(tUserEmail); err != nil {
		t.Fatal(err)
	}

	if sentEmails := e.Sent(); len(sentEmails) == 0 {
		t.Errorf("ForgotPassword sent zero emails")
	}
}

func TestForgotPasswordCtrl_ForgotPassword_invalid_email(t *testing.T) {
	m := newMemberFinderTD()
	tk := newTokenGeneratorTD()
	e := newEmailerTD()
	ctrl := NewForgotPasswordCtrl(m, tk, &e)

	if err := ctrl.ForgotPassword("not.this.user@example.com"); err != nil {
		t.Fatal(err)
	}

	if sentEmails := e.Sent(); len(sentEmails) != 0 {
		t.Errorf("ForgotPassword sent %v emails, want zero", len(sentEmails))
		t.Log("sentEmails:", sentEmails)
	}
}

type MemberFinderTD struct {
	email string
}

func newMemberFinderTD() MemberFinderTD {
	return MemberFinderTD{}
}
func (td MemberFinderTD) ByEmail(em string) *Member {
	if em == tUserEmail {
		return NewMember(em, "")
	}
	return nil
}

type TokenGeneratorTD struct{}

func newTokenGeneratorTD() TokenGeneratorTD {
	return TokenGeneratorTD{}
}
func (td TokenGeneratorTD) GenerateResetToken() (*Token, error) {
	return &Token{}, nil
}

type EmailerTD struct {
	sent []string
}

func newEmailerTD() EmailerTD {
	return EmailerTD{}
}
func (td *EmailerTD) SendPasswordForgotten(email string, member *Member, token *Token) error {
	if member != nil {
		td.sent = append(td.sent, fmt.Sprintf("password-forgotten:%s:%s", token, email))
	}

	return nil
}
func (td EmailerTD) Sent() []string {
	return td.sent
}
