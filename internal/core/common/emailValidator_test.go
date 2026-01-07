package common

import (
	"testing"
)

func TestValidarEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "Email válido simples",
			email:    "user@example.com",
			expected: true,
		},
		{
			name:     "Email válido com subdomain",
			email:    "user@mail.example.com",
			expected: true,
		},
		{
			name:     "Email válido com números",
			email:    "user123@example.com",
			expected: true,
		},
		{
			name:     "Email válido com pontos",
			email:    "first.last@example.com",
			expected: true,
		},
		{
			name:     "Email válido com hífen",
			email:    "user-name@example.com",
			expected: true,
		},
		{
			name:     "Email válido com underscore",
			email:    "user_name@example.com",
			expected: true,
		},
		{
			name:     "Email válido com mais (+)",
			email:    "user+tag@example.com",
			expected: true,
		},
		{
			name:     "Email válido com percentual",
			email:    "user%test@example.com",
			expected: true,
		},
		{
			name:     "Email válido com domínio .br",
			email:    "user@example.br",
			expected: true,
		},
		{
			name:     "Email válido com domínio .info",
			email:    "user@example.info",
			expected: true,
		},
		{
			name:     "Email inválido - sem @",
			email:    "userexample.com",
			expected: false,
		},
		{
			name:     "Email inválido - sem domínio",
			email:    "user@",
			expected: false,
		},
		{
			name:     "Email inválido - sem usuário",
			email:    "@example.com",
			expected: false,
		},
		{
			name:     "Email inválido - sem extensão",
			email:    "user@example",
			expected: false,
		},
		{
			name:     "Email inválido - extensão muito curta",
			email:    "user@example.c",
			expected: false,
		},
		{
			name:     "Email inválido - extensão muito longa",
			email:    "user@example.commmm",
			expected: false,
		},
		{
			name:     "Email inválido - com espaços",
			email:    "user name@example.com",
			expected: false,
		},
		{
			name:     "Email inválido - com maiúsculas",
			email:    "User@Example.COM",
			expected: false,
		},
		{
			name:     "Email inválido - dois @",
			email:    "user@@example.com",
			expected: false,
		},
		{
			name:     "Email vazio",
			email:    "",
			expected: false,
		},
		{
			name:     "Email inválido - ponto no início",
			email:    ".user@example.com",
			expected: false,
		},
		{
			name:     "Email inválido - ponto no final antes do @",
			email:    "user.@example.com",
			expected: false,
		},
		{
			name:     "Email inválido - caracteres especiais inválidos",
			email:    "user#name@example.com",
			expected: false,
		},
		{
			name:     "Email inválido - domínio com espaço",
			email:    "user@exam ple.com",
			expected: false,
		},
		{
			name:     "Email válido com múltiplos pontos no domínio",
			email:    "user@mail.server.example.com",
			expected: true,
		},
		{
			name:     "Email válido com números no domínio",
			email:    "user@example123.com",
			expected: true,
		},
		{
			name:     "Email válido com hífen no domínio",
			email:    "user@ex-ample.com",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidarEmail(tt.email)
			if result != tt.expected {
				t.Errorf("ValidarEmail(%s) = %v, esperado %v", tt.email, result, tt.expected)
			}
		})
	}
}

func TestValidarEmail_EdgeCases(t *testing.T) {
	t.Run("Email com todos caracteres permitidos", func(t *testing.T) {
		result := ValidarEmail("user.name+tag%123_test-email@sub.domain-name.example.com")
		expected := true
		if result != expected {
			t.Errorf("Falhou ao validar email complexo válido")
		}
	})

	t.Run("Email com extensão de 2 caracteres", func(t *testing.T) {
		result := ValidarEmail("user@example.co")
		expected := true
		if result != expected {
			t.Errorf("Falhou ao validar email com extensão de 2 caracteres")
		}
	})

	t.Run("Email com extensão de 4 caracteres", func(t *testing.T) {
		result := ValidarEmail("user@example.info")
		expected := true
		if result != expected {
			t.Errorf("Falhou ao validar email com extensão de 4 caracteres")
		}
	})
}

func BenchmarkValidarEmail(b *testing.B) {
	email := "user@example.com"
	for i := 0; i < b.N; i++ {
		ValidarEmail(email)
	}
}

func BenchmarkValidarEmail_Complex(b *testing.B) {
	email := "user.name+tag@sub.domain.example.com"
	for i := 0; i < b.N; i++ {
		ValidarEmail(email)
	}
}

func BenchmarkValidarEmail_Invalid(b *testing.B) {
	email := "invalid-email"
	for i := 0; i < b.N; i++ {
		ValidarEmail(email)
	}
}
