package iammember

import "strings"

// Parse a member to extract the kind and value.
func Parse(member string) (kind, value string, ok bool) {
	indexOfColon := strings.IndexByte(member, ':')
	if indexOfColon == -1 || indexOfColon == 0 || indexOfColon == len(member)-1 {
		return "", "", false
	}
	return member[:indexOfColon], member[indexOfColon+1:], true
}
