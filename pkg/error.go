package pkg

import "log"

// Generic Error
type GenericError struct{ message string }

func (b GenericError) Error() string {
	return "ERROR: " + b.message
}

// Check
func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
