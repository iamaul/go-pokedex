package csrf

import (
	"crypto/sha256"
	"encoding/base64"
	"io"

	"github.com/iamaul/go-pokedex/pkg/logger"
)

const (
	CSRFHeader = "X-CSRF-Token"
	csrfSalt   = "KbWaoi5xtDC3GEfBa9ovQdzOzXsuVU9I"
)

func MakeToken(sid string, logger logger.Logger) string {
	hash := sha256.New()
	_, err := io.WriteString(hash, csrfSalt+sid)
	if err != nil {
		logger.Errorf("Make CSRF Token", err)
	}
	token := base64.RawStdEncoding.EncodeToString(hash.Sum(nil))
	return token
}

func ValidateToken(token string, sid string, logger logger.Logger) bool {
	trueToken := MakeToken(sid, logger)
	return token == trueToken
}
