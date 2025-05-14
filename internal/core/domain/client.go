package domain

import (
	"errors"
	"regexp"
	"time"
)

type Cliente struct {
	ID        string    `json:"id"`
	Nome      string    `json:"nome"`
	CPF       string    `json:"cpf"`
	Email     string    `json:"email"`
	Telefone  string    `json:"telefone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NovoCliente(id, nome, cpf, email, telefone string) (*Cliente, error) {
	cliente := &Cliente{
		ID:        id,
		Nome:      nome,
		CPF:       cpf,
		Email:     email,
		Telefone:  telefone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := cliente.Validar()
	if err != nil {
		return nil, err
	}

	return cliente, nil
}

func (c *Cliente) Validar() error {
	if c.Nome == "" {
		return errors.New("nome não pode ser vazio")
	}

	if !ValidarCPF(c.CPF) {
		return errors.New("CPF inválido")
	}

	if !ValidarEmail(c.Email) {
		return errors.New("email inválido")
	}

	if c.Telefone == "" {
		return errors.New("telefone não pode ser vazio")
	}

	return nil
}

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

func ValidarEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return re.MatchString(email)
}
