package services

import (
	"regexp"
	"strings"
)

var dniArgentino7u8Digitos = regexp.MustCompile(`^\d{7,8}$`)

/* ValidateDNIAr exige 7 u 8 digitos numericos solamente (sin puntos ni letras). */
func ValidateDNIAr(dni string) (msg string, ok bool) {
	d := strings.TrimSpace(dni)
	if d == "" {
		return "El DNI es obligatorio", false
	}
	if !dniArgentino7u8Digitos.MatchString(d) {
		return "El DNI debe tener 7 u 8 digitos numericos (sin puntos ni letras)", false
	}
	return "", true
}
