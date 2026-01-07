package common

import (
	"testing"
)

func TestValidarCPF(t *testing.T) {
	tests := []struct {
		name     string
		cpf      string
		expected bool
	}{
		{
			name:     "CPF válido com pontuação",
			cpf:      "123.456.789-09",
			expected: true,
		},
		{
			name:     "CPF válido sem pontuação",
			cpf:      "12345678909",
			expected: true,
		},
		{
			name:     "CPF válido - caso real 1",
			cpf:      "111.444.777-35",
			expected: true,
		},
		{
			name:     "CPF válido - caso real 2",
			cpf:      "52998224725",
			expected: true,
		},
		{
			name:     "CPF inválido - todos os dígitos iguais (11111111111)",
			cpf:      "111.111.111-11",
			expected: false,
		},
		{
			name:     "CPF inválido - todos os dígitos iguais (00000000000)",
			cpf:      "000.000.000-00",
			expected: false,
		},
		{
			name:     "CPF inválido - todos os dígitos iguais (99999999999)",
			cpf:      "999.999.999-99",
			expected: false,
		},
		{
			name:     "CPF inválido - dígito verificador errado",
			cpf:      "123.456.789-10",
			expected: false,
		},
		{
			name:     "CPF inválido - tamanho incorreto (muito curto)",
			cpf:      "123.456.789",
			expected: false,
		},
		{
			name:     "CPF inválido - tamanho incorreto (muito longo)",
			cpf:      "123.456.789-099",
			expected: false,
		},
		{
			name:     "CPF vazio",
			cpf:      "",
			expected: false,
		},
		{
			name:     "CPF com letras",
			cpf:      "123.456.789-0a",
			expected: false,
		},
		{
			name:     "CPF com caracteres especiais",
			cpf:      "123@456#789-09",
			expected: false,
		},
		{
			name:     "CPF válido com espaços",
			cpf:      "123 456 789 09",
			expected: true,
		},
		{
			name:     "CPF inválido - primeiro dígito verificador incorreto",
			cpf:      "12345678808",
			expected: false,
		},
		{
			name:     "CPF inválido - segundo dígito verificador incorreto",
			cpf:      "12345678900",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidarCPF(tt.cpf)
			if result != tt.expected {
				t.Errorf("ValidarCPF(%s) = %v, esperado %v", tt.cpf, result, tt.expected)
			}
		})
	}
}

func TestValidarCPF_EdgeCases(t *testing.T) {
	t.Run("CPF com apenas números zero exceto dígitos verificadores", func(t *testing.T) {
		result := ValidarCPF("00000000191")
		expected := true
		if result != expected {
			t.Errorf("ValidarCPF(00000000191) = %v, esperado %v", result, expected)
		}
	})

	t.Run("CPF com múltiplos tipos de separadores", func(t *testing.T) {
		result := ValidarCPF("123.456.789-09")
		expected := true
		if result != expected {
			t.Errorf("Falhou ao validar CPF com separadores padrão")
		}
	})
}

func BenchmarkValidarCPF(b *testing.B) {
	cpf := "123.456.789-09"
	for i := 0; i < b.N; i++ {
		ValidarCPF(cpf)
	}
}

func BenchmarkValidarCPF_Invalid(b *testing.B) {
	cpf := "111.111.111-11"
	for i := 0; i < b.N; i++ {
		ValidarCPF(cpf)
	}
}
