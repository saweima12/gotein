package utils

import "strings"

func IsCommand(prefix string, text string) bool {
	return strings.HasPrefix(text, prefix)
}

func GetCommand(prefix string, text string) (cmd string, argStr string, ok bool) {
	if !strings.HasPrefix(text, prefix) {
		return "", "", false
	}

	shards := strings.Split(text, " ")

	if len(shards) > 0 {
		cmd = shards[0][1:]
	}

	if len(shards) == 1 {
		return cmd, "", true
	}

	// shards >= 2
	return cmd, shards[1], true
}
