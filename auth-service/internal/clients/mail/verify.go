package mail

import (
	emailverifier "github.com/AfterShip/email-verifier"
)

var verifier = emailverifier.NewVerifier()

func VerifyMail(email string) error {
	_, err := verifier.Verify(email)
	return err
}
