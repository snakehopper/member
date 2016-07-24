package member

import (
	"testing"
)

func TestCreditCtrl_CurrentCredit(t *testing.T) {
	m := newCreditUserFinderTD(tUserEmail)
	w := newCreditWalleterTD()
	ctrl := NewCreditCtrl(m, &w)

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
	ctrl := NewCreditCtrl(m, &w)

	if _, err := ctrl.CurrentCredit(); err != ErrEmptyMember {
		t.Errorf("CurrentCredit return error %v, want %v", err, ErrEmptyMember)
	}
}

func TestCreditCtrl_AddCredit(t *testing.T) {
	m := newCreditUserFinderTD(tUserEmail)
	w := newCreditWalleterTD()
	ctrl := NewCreditCtrl(m, &w)

	current, _ := ctrl.CurrentCredit()
	if total, err := ctrl.AddCredit(tCredit); err != nil {
		t.Fatal(err)
	} else if total != current+tCredit {
		t.Errorf("IncreaseCredit return %v, want %v", total, current+tCredit)
	}

	if current2, _ := ctrl.CurrentCredit(); current2 != current+tCredit {
		t.Errorf("CurrentCredit return %v, want %v", current2, current+tCredit)
	}
}

type CreditWalleterTD struct {
	amount float64
}

func newCreditWalleterTD() CreditWalleterTD {
	return CreditWalleterTD{amount: tCredit}
}

func (w CreditWalleterTD) ByUserId(id string) (float64, error) {
	if id == tUserEmail {
		return w.amount, nil
	}
	return 0, ErrEmptyMember
}

func (w *CreditWalleterTD) AddCredit(uid string, v float64) (float64, error) {
	w.amount += v
	return w.amount, nil
}

func (w CreditWalleterTD) OnAddCreditAlert(string, Wallet) {}

func (w CreditWalleterTD) String() string {
	return "CreditWalleterTD"
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
