package infrastructure

import "golang.org/x/crypto/bcrypt"

type BcryptPassword struct{ cost int }

func NewBcryptPassword(cost int) *BcryptPassword { return &BcryptPassword{cost: cost} }

func (b *BcryptPassword) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	return string(hash), err
}

func (b *BcryptPassword) Verify(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
