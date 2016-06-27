package member

type CreditCtrl struct {
	user   CreditUserFinder
	credit CreditWalleter
}

type CreditUserFinder interface {
	GetUserId() string
}

type CreditWalleter interface {
	ByUserId(string) (float64, error)
}

func NewCreditCtrl(u CreditUserFinder, w CreditWalleter) *CreditCtrl {
	return &CreditCtrl{u, w}
}

func (c CreditCtrl) CurrentCredit() (float64, error) {
	mid := c.user.GetUserId()
	if mid == "" {
		return 0, ErrEmptyMember
	}

	return c.credit.ByUserId(mid)
}
