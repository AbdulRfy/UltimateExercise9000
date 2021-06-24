package main

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID       uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Name     string `json:"name"`
	Email    string `gorm:"size:100;not null;unique" json:"email"`
	Password string `gorm:"size:100;not null;" json:"password"`
}

func (u *User) VerifyPassword(password string) bool {
	err = VerifyPassword(u.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return false
	}
	return true
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
