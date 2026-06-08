package models

type Card struct {
	Number    string `json:"number"`
	Exp_Month int    `json:"exp_month"`
	Exp_Year  int    `json:"exp_year"`
	CVC       string `json:"cvc"`
	Name      string `json:"name"`
}

type Token struct {
	ID          string `db:"id" json:"id`
	MERCHANT_ID string `db:"merch_id" json:"merch_id"`
	VAULT_ID    string `db:"vault_id" json:"vault_id"`
	BRAND       string `db:"brand" json:"brand"`
	Last4       string `db:"last4" json:"last4"`
	ExpMonth    int16  `db:"exp_month" json:"exp_month"`
	ExpYear     int16  `db:"exp_year" json:"exp_year"`
	EXPIRES_AT  string `db:"expires_at" json:"expires_at"`
}

type TokenizeCardRequest struct {
	Card Card `json:"card"`
}
