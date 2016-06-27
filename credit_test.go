package member

import (
	"testing"
)

func TestCreditCtrl_CurrentCredit(t *testing.T) {
	m := newCreditUserFinderTD(tUserEmail)
	w := newCreditWalleterTD()
	ctrl := NewCreditCtrl(m, w)

	v, err := ctrl.CurrentCredit()
	if err != nil {
		t.Fatal(err)
	}

	if v != tCredit {
		t.Errorf("CurrentCredit return %v, want %v", v, tCredit)
	}
}

func TestCreditCtrl_CurrentCredit_nosuch_user(t *testing.T) {
	m := newCreditUserFinderTD("not.such@email.com")
	w := newCreditWalleterTD()
	ctrl := NewCreditCtrl(m, w)

	if _, err := ctrl.CurrentCredit(); err != ErrEmptyMember {
		t.Errorf("CurrentCredit return error %v, want %v", err, ErrEmptyMember)
	}
}

type CreditWalleterTD struct {
}

func newCreditWalleterTD() CreditWalleterTD {
	return CreditWalleterTD{}
}

func (w CreditWalleterTD) ByUserId(id string) (float64, error) {
	if id == tUserEmail {
		return tCredit, nil
	}
	return 0, ErrEmptyMember
}

type CreditUserFinderTD struct {
	email string
}

func newCreditUserFinderTD(em string) CreditUserFinderTD {
	return CreditUserFinderTD{email: em}
}

func (td CreditUserFinderTD) GetUserId() string {
	return td.email
}
