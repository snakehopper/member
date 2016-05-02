package member

type ForgotPasswordCtrl struct {
	member MemberFinder
	token  TokenGenerator
	email  Emailer
}

func NewForgotPasswordCtrl(m MemberFinder, t TokenGenerator, e Emailer) *ForgotPasswordCtrl {
	return &ForgotPasswordCtrl{m, t, e}
}

func (c ForgotPasswordCtrl) ForgotPassword(em string) error {
	u := c.member.ByEmail(em)
	if u == nil {
		//by returning a success, we don't disclose to a potential attacker
		//whether the email address belongs to a registered user or not
		return nil
	}

	token, err := c.token.GenerateResetToken()
	if err != nil {
		return err
	}

	return c.email.SendPasswordForgotten(em, u, token)
}
