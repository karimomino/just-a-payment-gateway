package tokens

import (
	"context"
	"database/sql"
	"just-a-payment-gateway/backend/internal/models"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	DB *sql.DB
}

func (e *Handler) SaveCardTransaciton(encryptedCard []byte, w http.ResponseWriter) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	tx, err := e.DB.BeginTx(ctx, nil)
	if err != nil {
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

	if err := tx.QueryRowContext(ctx, query,
		encryptedCard, "key_3").Scan(&vault_id); err != nil {
		return "", err
	}

	tx.Commit()

	return vault_id, nil
}

func (e *Handler) GenerateTokenTransaction(tokenInfo *models.Token, w http.ResponseWriter) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := e.DB.BeginTx(ctx, nil)
	if err != nil {
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
	err = tx.QueryRowContext(ctx, query,
		tokenInfo.MERCHANT_ID, tokenInfo.VAULT_ID, tokenInfo.BRAND, tokenInfo.Last4, tokenInfo.ExpMonth, tokenInfo.ExpYear, "active").Scan(&tokenInfo.ID, &tokenInfo.EXPIRES_AT)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}
