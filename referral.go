package member

import (
	"time"
)

type Referral struct {
	Code        string
	Name        string
	Credit      int
	CreatedDate time.Time
}

type ReferralRecord struct {
	Referrer    string
	Referred    string
	Code        string
	Credit      int
	CreatedDate time.Time
}

type ReferralCtrl struct {
	referral ReferralFinder
	member   ReferralMemberUpdater
}

type ReferralFinder interface {
	ByCode(string) *Referral
	InsertRecord(ReferralRecord) error
}

type ReferralMemberUpdater interface {
	MemberId() string
	IncreaseCredit(int) error
	PatchMemberState() error
	IsNewRegistered() bool
}

func NewReferralCtrl(r ReferralFinder, m ReferralMemberUpdater) *ReferralCtrl {
	return &ReferralCtrl{r, m}
}

func (c ReferralCtrl) ApplyReferralCode(code string) error {
	r := c.referral.ByCode(code)
	if r == nil {
		return ErrNoSuchReferralCode
	}

	if !c.member.IsNewRegistered() {
		return ErrReferRegisteredUser
	}

	rc := r.Credit
	if err := c.member.IncreaseCredit(rc); err != nil {
		return err
	}

	return c.member.PatchMemberState()
}

func (c ReferralCtrl) LogReferral(code string) error {
	r := c.referral.ByCode(code)
	if r == nil {
		return ErrNoSuchReferralCode
	}

	rc := r.Credit
	src := r.Name
	rcv := c.member.MemberId()
	record := ReferralRecord{src, rcv, code, rc, time.Now()}
	return c.referral.InsertRecord(record)
}
