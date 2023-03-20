package utilities

import (
	"errors"

	emailverifier "github.com/AfterShip/email-verifier"
)

func ValidateEmail(email string) error {
	verifier := emailverifier.NewVerifier()

	res, err := verifier.Verify(email)
	if err != nil || !res.Syntax.Valid {
		return errors.New("invalid email address")
	}

	return nil
}
