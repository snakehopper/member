package member

type Member struct {
	Email string
	Name  string
}

func NewMember(em, name string) *Member {
	return &Member{Email: em, Name: name}
}

type MemberFinder interface {
	ByEmail(string) *Member
}
