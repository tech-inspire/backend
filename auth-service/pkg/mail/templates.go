package mail

import "fmt"

type Message struct {
	Title   string
	Message string
}

func ConfirmEmail(code string) Message {
	return Message{
		Title:   "Inspire: Confirm your email",
		Message: fmt.Sprintf("Your confirmation code is %s", code),
	}
}

func ResetPassword(code string) Message {
	const subject = "Inspire: Reset Password"

	return Message{
		Title:   "Inspire: Reset Password",
		Message: fmt.Sprintf("Your password reset code is %s", code),
	}
}
