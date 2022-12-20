package util

import (
	"fmt"
	"log"
)

func LogDebug(text string) {
	log.Print(text)
	fmt.Print(text)
}

func LogError(errorText string) {
	log.Fatal(errorText)
	fmt.Print(errorText)
}
