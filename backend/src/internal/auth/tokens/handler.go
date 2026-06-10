package tokens

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"just-a-payment-gateway/backend/internal/crypto"
	"just-a-payment-gateway/backend/internal/models"
)

type CardBrandSpace struct {
	VISA       int // 400000 - 499999 pan length 16
	MASTERCARD int // 222100 - 272099 and 510000 - 559999 pan length 16
	AMEX       int // 340000 - 349999 and 370000 - 379999 pan length 15
	UNKOWN     int // couldnt identify
	ERROR      int
}

var CardBrand = CardBrandSpace{
	VISA:       1,
	MASTERCARD: 2,
	AMEX:       3,
	UNKOWN:     4,
	ERROR:      -1,
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

func bindRequest(r *http.Request, card *models.Card) error {
	var request models.TokenizeCardRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Fatal(err)
		return err
	}
	*card = request.Card
	replacer := strings.NewReplacer(" ", "", "-", "")
	card.Number = replacer.Replace(card.Number)
	return nil
}

type ApiError struct {
	StatusCode int    `json:"status"`
	Error      string `json:"error,omitempty"`
}

func SendError(statusCode int, message string, w http.ResponseWriter) {
	error := ApiError{
		StatusCode: statusCode,
		Error:      message,
	}

	w.WriteHeader(error.StatusCode)
	json.NewEncoder(w).Encode(error)
}

func SendSuccess(statusCode int, data any, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (e *Handler) PostTokenizeCard(w http.ResponseWriter, r *http.Request) {
	var card models.Card
	w.Header().Set("Content-Type", "application/json")

	if err := bindRequest(r, &card); err != nil {
		SendError(http.StatusPaymentRequired, err.Error(), w)
		return
	}

	if err := runChecks(card); err != nil {
		SendError(http.StatusPaymentRequired, err.Error(), w)
		return
	}

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)

	if err := enc.Encode(card); err != nil {
		panic(err)
	}
	byteData := buff.Bytes()
	encrypted_card, _, err := crypto.EncryptPAN([]byte(byteData))
	if err != nil {
		SendError(http.StatusInternalServerError, "", w)
		return
	}

	vault_id, err := e.SaveCardTransaciton(encrypted_card, w)
	if err != nil {
		SendError(http.StatusPaymentRequired, err.Error(), w)
		return
	}

	token := models.Token{
		VAULT_ID:    vault_id,
		MERCHANT_ID: "merch_1",
		BRAND:       resolveBrand(GetBrand(card.Number[:6])),
		Last4:       card.Number[len(card.Number)-4:],
		ExpMonth:    int16(card.Exp_Month),
		ExpYear:     int16(card.Exp_Year),
	}

	if err := e.GenerateTokenTransaction(&token, w); err != nil {
		SendError(http.StatusPaymentRequired, err.Error(), w)
		return
	}

	SendSuccess(http.StatusCreated, token, w)
}

func runChecks(card models.Card) error {
	if !DateChecksPassed(card.Exp_Month, card.Exp_Year) {
		return errors.New("Card is Expired.")
	}

	passed, err := LuhnCheckPassed(card.Number)
	if passed != true {
		return err
	}

	return nil
}
