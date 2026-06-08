package tokens

import "strconv"

func LuhnCheckPassed() bool {
	return true
}

func DateChecksPassed() bool {
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
