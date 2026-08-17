package bcrypt

import "golang.org/x/crypto/bcrypt"

type Hasher struct{}

func NewHasher() *Hasher {
	return &Hasher{}
}

func (h *Hasher) Hash(password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(result), err
}

func (h *Hasher) Compare(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
