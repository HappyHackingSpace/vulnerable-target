package utils

import "fmt"

var rainbowColors = []string{
	"\033[31m", // Red
	"\033[33m", // Yellow
	"\033[32m", // Green
	"\033[36m", // Cyan
	"\033[34m", // Blue
	"\033[35m", // Magenta
}

func RainbowText(text string) string {
	runes := []rune(text)
	output := ""
	for i, r := range runes {
		color := rainbowColors[i%len(rainbowColors)]
		output += fmt.Sprintf("%s%c", color, r)
	}
	output += "\033[0m" // Reset
	return output
}
