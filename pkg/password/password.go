package password

type Password interface {
	Create(password string) (string, error)
	Compare(hashedPassword string, plainPassword string) error
}
