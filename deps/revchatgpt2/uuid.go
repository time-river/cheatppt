package revchatgpt2

import (
	"regexp"

	"github.com/google/uuid"
)

func isValidUUIDv4(str string) bool {
	uuidv4Re := regexp.MustCompile("(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$")
	return len(str) > 0 && uuidv4Re.MatchString(str)
}

func uuidv4() string {
	return uuid.Must(uuid.NewRandom()).String()
}
