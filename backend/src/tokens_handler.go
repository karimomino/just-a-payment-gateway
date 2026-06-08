package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

func (e *EnvDB) SaveCardTransaciton(encryptedCard []byte, c *gin.Context) (string, error) {
	tx, err := e.DB.BeginTx(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction: " + err.Error()})
		return "", err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Transaction rollback failed: %v", rollbackErr)
			}
		}
	}()

	query := "INSERT INTO vault (encrypted_card, encryption_key_id ) VALUES ($1, $2) RETURNING vault_id"
	var vault_id string
	err = tx.QueryRowContext(c.Request.Context(), query,
		encryptedCard, "key_3").Scan(&vault_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong."})
		return "", err
	}

	tx.Commit()

	return vault_id, nil
}

// func (e *EnvDB) SaveCardTransaciton(encryptedCard EncryptedCard, c *gin.Context) {
// 	tx, err := e.DB.BeginTx(c.Request.Context(), nil)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction: " + err.Error()})
// 		return
// 	}

// 	defer func() {
// 		if err != nil {
// 			if rollbackErr := tx.Rollback(); rollbackErr != nil {
// 				log.Printf("Transaction rollback failed: %v", rollbackErr)
// 			}
// 		}
// 	}()

// 	_, err = tx.ExecContext(c.Request.Context(),
// 		"INSERT INTO vault (encrypted_card, iv, pan_hash, bin, last4, exp_month, exp_year) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
// 		encryptedCard.EncryptedPAN, encryptedCard.IV, encryptedCard.PANHash, encryptedCard.BIN, encryptedCard.Last4, encryptedCard.ExpMonth, encryptedCard.ExpYear)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong."})
// 		return
// 	}

// 	tx.Commit()
// }

func (e *EnvDB) GenerateTokenTransaction(tokenInfo *Token, c *gin.Context) error {
	tx, err := e.DB.BeginTx(c.Request.Context(), nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction: " + err.Error()})
		return err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Transaction rollback failed: %v", rollbackErr)
			}
		}
	}()
	query := "INSERT INTO tokens (merch_id, vault_id, brand, last4, exp_month, exp_year, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, expires_at"
	err = tx.QueryRowContext(c.Request.Context(), query,
		tokenInfo.MERCHANT_ID, tokenInfo.VAULT_ID, tokenInfo.BRAND, tokenInfo.Last4, tokenInfo.ExpMonth, tokenInfo.ExpYear, "active").Scan(&tokenInfo.ID, &tokenInfo.EXPIRES_AT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong." + err.Error()})
		return err
	}

	tx.Commit()

	return nil
}

func (e *EnvDB) postTokenizeCard(c *gin.Context) {
	var newCardRequest TokenizeCardRequest

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

	encrypted_card, _, err := encryptPAN([]byte(byteData))
	if err != nil {
		return
	}

	vault_id, err := e.SaveCardTransaciton(encrypted_card, c)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	log.Println("Saved card to vault")

	token := Token{
		VAULT_ID:    vault_id,
		MERCHANT_ID: "merch_1",
		BRAND:       resolveBrand(GetBrand(card.Number[:6])),
		Last4:       card.Number[len(card.Number)-4:],
		ExpMonth:    int16(card.Exp_Month),
		ExpYear:     int16(card.Exp_Year),
	}
	log.Println("Created token variable")
	e.GenerateTokenTransaction(&token, c)

	c.IndentedJSON(http.StatusCreated, token)
}

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
