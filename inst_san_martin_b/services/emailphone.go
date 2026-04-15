package services

import (
	"regexp"
	"strings"
)

var (
	correoElectronicoRegex = regexp.MustCompile(`(?i)^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	telefonoSoloDigitos    = regexp.MustCompile(`^\d+$`)
)

/* ValidateCorreoElectronico: correo vacío = válido (campo opcional). Si hay texto, debe ser un correo con @ y dominio con extensión. */
func ValidateCorreoElectronico(email string) (msg string, ok bool) {
	e := strings.TrimSpace(email)
	if e == "" {
		return "", true
	}
	if len(e) > 254 {
		return "El correo electronico es demasiado largo", false
	}
	if !correoElectronicoRegex.MatchString(e) {
		return "El correo electronico debe ser valido (incluir @ y dominio con extension, sin espacios)", false
	}
	return "", true
}

/* ValidateCorreoElectronicoRequired exige correo no vacío y formato válido */
func ValidateCorreoElectronicoRequired(email string) (msg string, ok bool) {
	e := strings.TrimSpace(email)
	if e == "" {
		return "El correo electronico es obligatorio", false
	}
	return ValidateCorreoElectronico(e)
}

/* NormalizeTelefonoDigits deja solo dígitos (útil tras pegar números con espacios o prefijos). */
func NormalizeTelefonoDigits(phone string) string {
	return regexp.MustCompile(`\D`).ReplaceAllString(strings.TrimSpace(phone), "")
}

/* ValidateTelefono: vacío = válido. Con dígitos: entre 7 y 15 (nacional / internacional). */
func ValidateTelefono(digitsOnly string) (msg string, ok bool) {
	if digitsOnly == "" {
		return "", true
	}
	if !telefonoSoloDigitos.MatchString(digitsOnly) {
		return "El telefono solo puede contener numeros", false
	}
	if len(digitsOnly) < 7 || len(digitsOnly) > 15 {
		return "El telefono debe tener entre 7 y 15 digitos", false
	}
	return "", true
}
