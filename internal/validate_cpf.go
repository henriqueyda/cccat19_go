package internal

import (
	"fmt"
	"regexp"
)

var cleanCpfRegex = regexp.MustCompile(`\D`)

func validateCpf(cpf string) bool {
	if cpf == "" {
		return false
	}
	cpf = clean(cpf)
	if len(cpf) != 11 {
		return false
	}
	if allDigitsTheSame(cpf) {
		return false
	}
	firstDigit := calculateDigit(cpf, 10)
	secondDigit := calculateDigit(cpf, 11)
	actualDigit := extractDigit(cpf)
	return actualDigit == fmt.Sprintf("%d%d", firstDigit, secondDigit)
}

func clean(cpf string) string {
	return cleanCpfRegex.ReplaceAllString(cpf, "")
}

func allDigitsTheSame(cpf string) bool {
	firstDigit := cpf[0:1]
	for _, digit := range cpf {
		if string(digit) != firstDigit {
			return false
		}
	}
	return true
}

func calculateDigit(cpf string, factor int) int {
	total := 0
	for _, digit := range cpf {
		if factor > 1 {
			total += int(digit-'0') * factor
			factor--
		}
	}
	remainder := total % 11
	if remainder < 2 {
		return 0
	}
	return 11 - remainder
}

func extractDigit(cpf string) string {
	return cpf[9:]
}
