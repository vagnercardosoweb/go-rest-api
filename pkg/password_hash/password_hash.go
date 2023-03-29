package password_hash

type PasswordHash interface {
	Create(password string) (string, error)
	Compare(hashedPassword string, plainPassword string) error
}
