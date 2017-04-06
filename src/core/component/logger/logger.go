package app_logger

import "log"

func Error(message string) {
	log.Println("Error: ", message)
}
