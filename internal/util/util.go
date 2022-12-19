package util

import "fmt"

func LogDebug(text string) {
	fmt.Print(text)
}

func LogError(errorText string) {
	fmt.Print(errorText)
}
