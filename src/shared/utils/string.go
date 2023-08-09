package utils

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func GenerateSlug(text string) string {
	text = strings.ToLower(text)
	reg := regexp.MustCompile("[^a-zA-Z0-9\\s]+")

	text = reg.ReplaceAllString(text, "")
	text = strings.ReplaceAll(text, " ", "-")
	text = regexp.MustCompile("-+").ReplaceAllString(text, "-")
	text = strings.Trim(text, "-") + "-" + time.Now().UTC().Format("20060102150405")

	return text
}

func GetRandomString(data []string) string {
	idx := rand.Intn(len(data))
	return data[idx]
}
