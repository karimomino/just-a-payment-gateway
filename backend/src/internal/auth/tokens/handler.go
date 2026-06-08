package tokens

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"

	"just-a-payment-gateway/backend/internal/crypto"
	"just-a-payment-gateway/backend/internal/models"

	"github.com/gin-gonic/gin"
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

func (e *Handler) PostTokenizeCard(c *gin.Context) {
	var newCardRequest models.TokenizeCardRequest

	if err := c.BindJSON(&newCardRequest); err != nil {
		log.Fatal(err)
		return
	}
	card := newCardRequest.Card

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)

	if err := enc.Encode(card); err != nil {
		panic(err)
	}
	byteData := buff.Bytes()

	fmt.Println(byteData)
	encrypted_card, _, err := crypto.EncryptPAN([]byte(byteData))
	if err != nil {
		return
	}

	vault_id, err := e.SaveCardTransaciton(encrypted_card, c)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
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

	e.GenerateTokenTransaction(&token, c)

	c.IndentedJSON(http.StatusCreated, token)
}
