package helper

import "golang.org/x/crypto/bcrypt"

const cost = 12

func Hash(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), cost)
	return string(b), err
}

func Verify(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
