package utils

import "log"


func SendEmailWarning(email string) {
	log.Printf("Warning: IP address changed for user %s", email)
}
