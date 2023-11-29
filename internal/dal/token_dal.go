package dal

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"encoding/hex"
	"fanclub/internal/models"
	"time"
)

type TokenDAL interface {
	RefreshUserToken(userId int, ttl time.Duration) (string, error)
	GetUserFromToken(token string) (int, bool, error)
	Insert(model models.Token) error
}

type tokenDAL struct {
	DB *sql.DB
}

var _ TokenDAL = (*tokenDAL)(nil)

func NewTokenDAL(db *sql.DB) TokenDAL {
	return &tokenDAL{db}
}

func (dal *tokenDAL) RefreshUserToken(userId int, ttl time.Duration) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	data := models.Token{
		UserID: userId,
		Hash:   token,
		Expiry: time.Now().Add(ttl),
	}

	if err := dal.Insert(data); err != nil {
		return "", err
	}
	return token, nil
}

func (dal *tokenDAL) Insert(model models.Token) error {
	query := `
	REPLACE INTO tokens(user_id, hash, expiry)
	VALUES(?, ?, ?);
	`

	_, err := dal.DB.Exec(query, model.UserID, model.Hash, model.Expiry)
	return err
}

// Return value: user_id, is_expired, error
func (dal *tokenDAL) GetUserFromToken(token string) (int, bool, error) {
	query := `SELECT user_id, expiry FROM tokens WHERE hash = ?;`
	var userID int
	var expiry time.Time
	if err := dal.DB.QueryRow(query, token).Scan(&userID, &expiry); err != nil {
		return 0, true, err
	}

	isExpired := time.Now().After(expiry)
	return userID, isExpired, nil
}

func generateToken() (string, error) {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	plaintext := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	sum256 := sha256.Sum256([]byte(plaintext))
	return hex.EncodeToString(sum256[:]), nil
}
