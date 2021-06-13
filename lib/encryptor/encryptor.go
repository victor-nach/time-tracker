package encryptor

import (
	"golang.org/x/crypto/bcrypt"
)

type Encryptor interface {
	ComparePasscode(passcode, hashedPasscode string) bool
	HashPassword(password string) (string, error)
}

type encryptor struct{}

func NewEncryptor() Encryptor {
	return &encryptor{}
}

func (e *encryptor) ComparePasscode(passcode, hashedPasscode string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPasscode), []byte(passcode))
	return err == nil
}

func (e *encryptor) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
