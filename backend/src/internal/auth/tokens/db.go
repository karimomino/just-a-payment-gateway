package tokens

import (
	"database/sql"
	"just-a-payment-gateway/backend/internal/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	DB *sql.DB
}

func (e *Handler) SaveCardTransaciton(encryptedCard []byte, c *gin.Context) (string, error) {
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

func (e *Handler) GenerateTokenTransaction(tokenInfo *models.Token, c *gin.Context) error {
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
