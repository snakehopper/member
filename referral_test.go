package member

import (
	"testing"
)

const (
	tReferralCode = "SOME_TEXT"
	tCredit       = 10
)

func TestReferralCtrl_ApplyReferralCode(t *testing.T) {
	r := newReferralFinderTD()
	m := newMemberUpdaterTD()
	ctrl := NewReferralCtrl(&r, &m)

	if !m.IsNewRegistered() {
		t.Errorf("Member should be a new registered state")
	}

	if err := ctrl.ApplyReferralCode(tReferralCode); err != nil {
		t.Fatal(err)
	}

	if m.credit != tCredit {
		t.Errorf("ApplyReferralCode credit not match")
	}

	if m.IsNewRegistered() {
		t.Errorf("ApplyReferralCode should update user registered state")
	}
}

func TestReferralCtrl_ApplyReferralCode_invalid_code(t *testing.T) {
	r := newReferralFinderTD()
	m := newMemberUpdaterTD()
	ctrl := NewReferralCtrl(&r, &m)

	if err := ctrl.ApplyReferralCode("not.such.code"); err != ErrNoSuchReferralCode {
		t.Fatalf("ApplyReferralCode should return error %v", err)
	}
}

func TestReferralCtrl_ApplyReferralCode_not_new_registered_user(t *testing.T) {
	r := newReferralFinderTD()
	m := MemberUpdaterTD{tUserEmail, tCredit, true}
	ctrl := NewReferralCtrl(&r, &m)

	if err := ctrl.ApplyReferralCode(tReferralCode); err != ErrReferRegisteredUser {
		t.Fatalf("ApplyReferralCode should return error %v", err)
	}
}

func TestReferralCtrl_LogReferral(t *testing.T) {
	r := newReferralFinderTD()
	m := newMemberUpdaterTD()
	ctrl := NewReferralCtrl(&r, &m)

	if err := ctrl.LogReferral(tReferralCode); err != nil {
		t.Fatal(err)
	}

	if len(r.records) == 0 {
		t.Errorf("LogReferral should insert log")
	}
}

type ReferralFinderTD struct {
	records []ReferralRecord
}

func newReferralFinderTD() ReferralFinderTD {
	return ReferralFinderTD{}
}
func (r ReferralFinderTD) ByCode(c string) *Referral {
	if c == tReferralCode {
		return &Referral{Code: c, Name: tUserEmail, Credit: tCredit}
	}
	return nil
}
func (r *ReferralFinderTD) InsertRecord(rc ReferralRecord) error {
	r.records = append(r.records, rc)
	return nil
}

type MemberUpdaterTD struct {
	id          string
	credit      int
	hasReferred bool
}

func newMemberUpdaterTD() MemberUpdaterTD {
	return MemberUpdaterTD{tUserEmail, tCredit, false}
}
func (td MemberUpdaterTD) MemberId() string { return td.id }
func (td MemberUpdaterTD) IncreaseCredit(cd int) error {
	td.credit += cd
	return nil
}
func (td *MemberUpdaterTD) PatchMemberState() error {
	td.hasReferred = true
	return nil
}
func (td MemberUpdaterTD) IsNewRegistered() bool { return !td.hasReferred }
