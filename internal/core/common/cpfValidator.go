package common

import (
	"regexp"
)

func ValidarCPF(cpf string) bool {

	re := regexp.MustCompile(`[^0-9]`)
	cpf = re.ReplaceAllString(cpf, "")

	if len(cpf) != 11 {
		return false
	}

	igual := true
	for i := 1; i < 11; i++ {
		if cpf[i] != cpf[0] {
			igual = false
			break
		}
	}
	if igual {
		return false
	}

	soma := 0
	for i := 0; i < 9; i++ {
		soma += int(cpf[i]-'0') * (10 - i)
	}
	resto := soma % 11
	if resto < 2 {
		if cpf[9] != '0' {
			return false
		}
	} else {
		if cpf[9] != byte('0'+11-resto) {
			return false
		}
	}

	soma = 0
	for i := 0; i < 10; i++ {
		soma += int(cpf[i]-'0') * (11 - i)
	}
	resto = soma % 11
	if resto < 2 {
		if cpf[10] != '0' {
			return false
		}
	} else {
		if cpf[10] != byte('0'+11-resto) {
			return false
		}
	}

	return true
}