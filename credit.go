package member

import (
	"math"
	"time"
)

type CreditCtrl struct {
	user   CreditUserFinder
	credit CreditWalleter
}

type CreditUserFinder interface {
	GetUserId() string
}

type CreditWalleter interface {
	String() string
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

func (c CreditCtrl) Origin() string {
	return c.credit.String()
}

type Wallet struct {
	Source     string
	Balance    float64
	Available  float64
	Processing float64 //Use to store amount to be process, and reset after MultiWalletCtrl.DoTransaction

	Error error
}
type MultiWalletCtrl struct {
	ctrls   []*CreditCtrl
	wallets map[string]Wallet
}

func NewMultiWalletCtrl(cc ...*CreditCtrl) MultiWalletCtrl {
	w := make([]*CreditCtrl, 0)
	w = append(w, cc...)
	return MultiWalletCtrl{ctrls: w}
}

func (w *MultiWalletCtrl) FetchCurrentCredit() {
	lg := len(w.ctrls)
	var ch = make(chan Wallet, lg)
	for i := 0; i < lg; i++ {
		go func(s *CreditCtrl) {
			v, err := s.CurrentCredit()
			ch <- Wallet{Source: s.Origin(), Balance: v, Available: v, Error: err}
		}(w.ctrls[i])
	}

	for i := 0; i < lg; i++ {
		src := w.ctrls[i].Origin()
		w.wallets[src] = <-ch
	}
}

func (w *MultiWalletCtrl) TradeCalculation(amount float64) []Wallet {
	var res = make([]Wallet, 0)

	for i := 0; i < len(w.ctrls); i++ {
		src := w.ctrls[i].Origin()
		wl := w.wallets[src]

		r := math.Min(amount, wl.Available)
		amount -= r
		wl.Available -= r
		wl.Processing -= r

		w.wallets[src] = wl
		res = append(res, wl)
	}

	return res
}

func (w MultiWalletCtrl) ProcessingAmount() float64 {
	var total float64
	for _, v := range w.wallets {
		total += v.Processing
	}
	return total
}

func (w *MultiWalletCtrl) DoTransaction() (errWallet []Wallet) {
	lg := len(w.ctrls)
	var ch = make(chan Wallet, lg)
	for i := 0; i < lg; i++ {
		go func(s *CreditCtrl) {
			tw := w.wallets[s.Origin()]
			for t := 0; t < 3; t++ {
				if _, tw.Error = s.AddCredit(tw.Processing); tw.Error == nil {
					tw.Processing = 0
					w.wallets[s.Origin()] = tw
					ch <- tw
					return
				}

				time.Sleep(300 * time.Millisecond)
			}
			w.wallets[s.Origin()] = tw
			ch <- tw
		}(w.ctrls[i])
	}

	errWallet = make([]Wallet, 0)
	for i := 0; i < lg; i++ {
		wl := <-ch
		if wl.Error != nil {
			errWallet = append(errWallet, wl)
		}
	}

	return errWallet
}
