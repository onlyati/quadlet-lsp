package semantic

import (
	"regexp"
)

// ImageValueTokens parses image value like: 'docker.io/gitea/gitea:rootless'
func ImageValueTokens(value string) []token {
	tokens := []token{}
	_ = regexp.MustCompile(`(?:[a-z0-9]+(?:[a-z0-9._-]+)*\.(?:[a-z0-9]+)|localhost)`)

	return tokens
}
