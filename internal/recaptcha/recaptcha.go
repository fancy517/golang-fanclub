package recaptcha

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type RecaptchaResponseType struct {
	Success            bool      `json:"success"`
	ChallengeTimestamp time.Time `json:"challenge_ts"`
	ErrorCodes         []string  `json:"error-codes"`
}

func VerifyRecaptcha(value string) (bool, error) {
	verifyUrl := os.Getenv("GOOGLE_RECAPTCHA_VERIFY_URL")
	secret := os.Getenv("GOOGLE_RECAPTCHA_SECRET_KEY")
	form := url.Values{}
	form.Set("secret", secret)
	form.Set("response", value)

	resp, err := http.Post(verifyUrl, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var r RecaptchaResponseType
	if err := json.Unmarshal(data, &r); err != nil {
		return false, err
	}

	return r.Success, nil
}
