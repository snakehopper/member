package member

import (
	"time"
)

type Referral struct {
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Credit      int       `json:"credit"`
	CreatedDate time.Time `json:"createdDate"`
}

type ReferralLog struct {
	Referrer    string    `json:"referrer"`
	Referred    string    `json:"referred"`
	Code        string    `json:"code"`
	Credit      int       `json:"credit"`
	CreatedDate time.Time `json:"createdDate"`
}

type ReferralCtrl struct {
	referral ReferralFinder
	member   ReferralMemberUpdater
}

type ReferralFinder interface {
	ByCode(string) *Referral
	InsertRecord(ReferralLog) error
}

type ReferralMemberUpdater interface {
	MemberId() string
	IncreaseCredit(int, string) error
	MarkedUserAsReferred(string) error
	IsValidForReferred() bool
}

func NewReferralCtrl(r ReferralFinder, m ReferralMemberUpdater) *ReferralCtrl {
	return &ReferralCtrl{r, m}
}

func (c ReferralCtrl) ApplyReferralCode(code string) error {
	r := c.referral.ByCode(code)
	if r == nil {
		return ErrNoSuchReferralCode
	}

	if !c.member.IsValidForReferred() {
		return ErrReferRegisteredUser
	}

	rc := r.Credit
	if err := c.member.IncreaseCredit(rc, code); err != nil {
		return err
	}

	return c.member.MarkedUserAsReferred(code)
}

func (c ReferralCtrl) LogReferral(code string) error {
	r := c.referral.ByCode(code)
	if r == nil {
		return ErrNoSuchReferralCode
	}

	rc := r.Credit
	src := r.Name
	rcv := c.member.MemberId()
	record := ReferralLog{src, rcv, code, rc, time.Now()}
	return c.referral.InsertRecord(record)
}
