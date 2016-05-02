package member

type Emailer interface {
	SendPasswordForgotten(email string, member *Member, token *Token) error
}
