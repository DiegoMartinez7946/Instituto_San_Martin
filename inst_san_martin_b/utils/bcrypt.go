package utils

import "golang.org/x/crypto/bcrypt"

/* Encrypt the uncoming password for save later */
func EncryptPassword(pass string) (string, error) {
	salt := 8

	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), salt)

	return string(bytes), err
}

/* Decrypt the password and compare with incoming password */
func DecryptPassword(password string, passwordIncoming string) (error) {
	passwordBytes := []byte(passwordIncoming)
	passwordBD := []byte(password)

	err := bcrypt.CompareHashAndPassword(passwordBD, passwordBytes)

	return err
}
