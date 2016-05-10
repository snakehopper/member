package member

import (
	"errors"
)

var (
	ErrNoSuchMember        = errors.New("No such member")
	ErrNoSuchReferralCode  = errors.New("No such referral code")
	ErrReferRegisteredUser = errors.New("Only new registered user can be referred")
)
