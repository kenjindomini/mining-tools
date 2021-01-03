package miningtools

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// checkError will log the error and execute additional handling (will not exit or panic on its own)
// any %s found in formattedMessage will be replaced with err.Error()
func checkError(err error, formattedMessage string, additionalErrorHandling func()) {
	if err != nil {
		fmt.Println(err)
		formattedMessage = strings.ReplaceAll(formattedMessage, "%s", err.Error())
		log.Error(formattedMessage)
		if additionalErrorHandling != nil {
			additionalErrorHandling()
		}
		return
	}
}

// checkPanic will execute addtional handling, if provided, then log.panic()
// any %s found in formattedMessage will be replaced with err.Error()
func checkPanic(err error, formattedMessage string, additionalErrorHandling func()) {
	if err != nil {
		fmt.Println(err)
		if additionalErrorHandling != nil {
			additionalErrorHandling()
		}
		formattedMessage = strings.ReplaceAll(formattedMessage, "%s", err.Error())
		log.Panic(formattedMessage)
		return
	}
}
