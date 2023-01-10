package util

import "net/mail"

func ValidMailAddress(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}
