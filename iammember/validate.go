package iammember

import (
	"fmt"
	"strings"
)

// Validate checks that an IAM member is valid.
func Validate(member string) error {
	if member == "" {
		return fmt.Errorf("member is empty")
	}
	indexOfColon := strings.IndexByte(member, ':')
	if indexOfColon == -1 || indexOfColon == 0 || indexOfColon == len(member)-1 {
		return fmt.Errorf("member '%s' has invalid format", member)
	}
	return nil
}
