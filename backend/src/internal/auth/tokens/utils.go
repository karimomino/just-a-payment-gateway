package tokens

import (
	"errors"
	"strconv"
	"time"
	"unicode"
)

func LuhnCheckPassed(cardNumber string) (bool, error) {
	cardNetwork := GetBrand(cardNumber)
	len := len(cardNumber)
	customError := errors.New("Invalid card number.")

	switch cardNetwork {
	case CardBrand.VISA:
		if len > 16 || len < 13 {
			return false, customError
		}
	case CardBrand.AMEX:
		if len != 15 {
			return false, customError
		}
	case CardBrand.MASTERCARD:
		if len != 16 {
			return false, customError
		}
	case CardBrand.ERROR, CardBrand.UNKOWN:
		return false, errors.New("Invalid car number or we don't support this car network.")
	default:
		break
	}

	if luhnCheck(cardNumber, len) != nil {
		return false, customError
	}

	return true, nil
}

func luhnCheck(cardNumber string, len int) error {
	sum := 0
	alternate := false

	for i := len - 1; i > -1; i-- {
		char := rune(cardNumber[i])

		if unicode.IsSpace(char) {
			continue
		}
		if !unicode.IsDigit(char) {
			return errors.New("Invalid card number.")
		}

		digit := int(char)

		if alternate {
			doubled := digit * 2
			if doubled > 19 {
				digit -= 9
			}
		}
		sum += digit
	}

	if sum%10 != 0 {
		return errors.New("Invalid card number.")
	}

	return nil
}

func DateChecksPassed(month int, year int) bool {
	now := time.Now()

	cardDate := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.Now().Location())

	if cardDate.Before(now) {
		return false
	}

	return true
}

func GetBrand(bin string) int {
	VISA_MIN := 400000
	VISA_MAX := 499999
	MASTERCARD_MIN := 222100
	MASTERCARD_MAX := 272099
	MASTERCARD_2_MIN := 510000
	MASTERCARD_2_MAX := 559999
	AMEX_MIN := 340000
	AMEX_MAX := 349999
	AMEX_2_MIN := 370000
	AMEX_2_MAX := 379999

	binI, err := strconv.Atoi(bin)
	if err != nil {
		return CardBrand.ERROR
	}

	if InRange(binI, VISA_MIN, VISA_MAX) {
		return CardBrand.VISA
	} else if InRange(binI, MASTERCARD_MIN, MASTERCARD_MAX) || InRange(binI, MASTERCARD_2_MIN, MASTERCARD_2_MAX) {
		return CardBrand.MASTERCARD
	} else if InRange(binI, AMEX_MIN, AMEX_MAX) || InRange(binI, AMEX_2_MIN, AMEX_2_MAX) {
		return CardBrand.AMEX
	}

	return CardBrand.UNKOWN
}

func resolveBrand(brand int) string {
	switch brand {
	case CardBrand.VISA:
		return "visa"
	case CardBrand.AMEX:
		return "amex"
	case CardBrand.MASTERCARD:
		return "mastercard"
	default:
		return "unknown"
	}
}

func InRange(bin int, min int, max int) bool {
	return bin >= min && bin <= max
}

func BrandCVCCheckPassed() bool {
	return true
}
