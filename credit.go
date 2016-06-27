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
	AddCredit(string, float64) (float64, error)
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

func (c CreditCtrl) AddCredit(amount float64) (float64, error) {
	uid := c.user.GetUserId()
	if uid == "" {
		return 0, ErrEmptyMember
	}

	return c.credit.AddCredit(uid, amount)
}
