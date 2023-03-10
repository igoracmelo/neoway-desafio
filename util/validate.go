package util

import (
	"fmt"
)

// `ValidateCpfOrCnpj` espera uma string `val` já sanitizada e verifica se
// `val` representa um CPF ou CNPJ matematicamente válido
func ValidateCpfOrCnpj(val string) error {
	if len(val) == 11 {
		return ValidateCpf(val)
	}
	if len(val) == 14 {
		return ValidateCnpj(val)
	}
	return fmt.Errorf("invalid CPF or CNPJ with length %d: %s", len(val), val)
}

// Retornar nil se o CPF for válido, senão retorna um erro especificando
// a validação que falhou
func ValidateCpf(val string) error {

	// checa se o CPF possui algum dígito distinto ou se então é formado só
	// por um dígito repetido, por exemplo: 555.555.555-55
	hasDistinctDigits := false
	for i := 1; i < 11; i++ {
		if val[i-1] != val[i] {
			hasDistinctDigits = true
			break
		}
	}

	if !hasDistinctDigits {
		return fmt.Errorf("Invalid CPF with no distinct digit: %s", val)
	}

	sum := 0

	// Somando os 9 primeiros dígitos do CPF, aplicando os seguintes pesos:
	// (A * 10) + (B * 9) + (C * 8) + (D * 7) + (E * 6) + (F * 5) + (G * 4) + (H * 3) + (I * 2)
	// onde as letras repesentariam um CPF ABC.DEF.GHI-JK
	for i := 0; i < 9; i++ {

		// converte o dígito em ascii para o número en inteiro
		// valores inválidos, como letras e símbolos não são tratados aqui,
		// mas sim na sanitização
		digit := int(val[i] - '0')
		sum += (10 - i) * digit
	}

	wantDigit10 := sum * 10 % 11 % 10 // dígito esperado
	gotDigit10 := int(val[9] - '0')   // dígito obtido

	if wantDigit10 != gotDigit10 {
		return fmt.Errorf("expected digit at position 10 to be %d, but got %d", wantDigit10, gotDigit10)
	}

	sum = 0

	// Somando os 10 primeiros dígitos do CPF, aplicando os seguintes pesos:
	// (A * 11) + (B * 10) + (C * 9) + (D * 8) + (E * 7) + (F * 6) + (G * 5) + (H * 4) + (I * 3) + (J * 2)
	// onde as letras repesentariam um CPF ABC.DEF.GHI-JK
	for i := 0; i < 10; i++ {
		digit := int(val[i] - '0')
		sum += (11 - i) * digit
	}

	wantDigit11 := sum * 10 % 11 % 10 // dígito esperado
	gotDigit11 := int(val[10] - '0')  // dígito obtido

	if wantDigit11 != gotDigit11 {
		return fmt.Errorf("expected digit at position 11 to be %d, but got %d", wantDigit11, gotDigit11)
	}

	return nil
}

func ValidateCnpj(val string) error {
	// checa se o CPF possui algum dígito distinto ou se então é formado só
	// por um dígito repetido, por exemplo: 555.555.555-55
	hasDistinctDigits := false
	for i := 1; i < 14; i++ {
		if val[i-1] != val[i] {
			hasDistinctDigits = true
			break
		}
	}

	if !hasDistinctDigits {
		return fmt.Errorf("Invalid CNPJ with no distinct digit: %s", val)
	}

	// cada dígito do CNPJ é multiplicado por um desses dígitos na etapa 1
	multipliers1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	// cada dígito do CNPJ é multiplicado por um desses dígitos na etapa 2
	multipliers2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	// similar ao cálculo do CPF
	sum := 0
	for i := 0; i < 12; i++ {
		digit := int(val[i] - '0')
		sum += digit * multipliers1[i]
	}

	// cálculo do primeiro dígito verificador
	wantDigit13 := sum % 11
	if wantDigit13 < 2 {
		wantDigit13 = 0
	} else {
		wantDigit13 = 11 - wantDigit13
	}

	gotDigit13 := int(val[12] - '0')

	if wantDigit13 != gotDigit13 {
		return fmt.Errorf("expected digit at position 13 to be %d, but got %d", wantDigit13, gotDigit13)
	}

	sum = 0
	for i := 0; i < 13; i++ {
		digit := int(val[i] - '0')
		sum += digit * multipliers2[i]
	}

	// cálculo do segundo dígito verificador
	wantDigit14 := sum % 11
	if wantDigit14 < 2 {
		wantDigit14 = 0
	} else {
		wantDigit14 = 11 - wantDigit14
	}

	gotDigit14 := int(val[13] - '0')

	if wantDigit14 != gotDigit14 {
		return fmt.Errorf("expected digit at position 14 to be %d, but got %d", wantDigit14, gotDigit14)
	}

	return nil
}
