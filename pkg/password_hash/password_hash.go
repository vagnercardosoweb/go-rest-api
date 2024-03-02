package password_hash

type PasswordHash interface {
	Compare(hashedPassword string, plainPassword string) error
	Create(password string) (string, error)
}
