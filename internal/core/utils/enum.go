package utils

import "strings"

type Action string

const (
	StatusOn  Action = "On"
	StatusOff Action = "Off"
)

func IsValidAction(at string) bool {
	switch Action(strings.ToLower(at)) {
	case Action(strings.ToLower(string(StatusOn))), Action(strings.ToLower(string(StatusOff))):
		return true
	}
	return false
}
