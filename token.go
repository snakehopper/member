package member

type TokenGenerator interface {
	GenerateResetToken() (*Token, error)
}

type Token struct {
	Raw string
}
