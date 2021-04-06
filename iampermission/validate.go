package iampermission

import (
	"fmt"
	"unicode"
)

// Validate a permission.
//
// A permission is valid on the format `<service>.<resource>.<verb>`.
//
// - It contains at most 3 segments.
// - When no segment is a wildcard, it contains exactly 3 segments.
// - All segments are non-empty.
// - All segments are lower-case.
// - The first (service) segment is not a wildcard.
// - A wildcard only occurs in the last segment.
func Validate(permission string) error {
	var sc Scanner
	sc.Init(permission)
	// Segment 1: Service.
	if !sc.Scan() {
		return fmt.Errorf("permission is empty")
	}
	if sc.Wildcard() {
		return fmt.Errorf("permission `%s`: service segment must not be wildcard", permission)
	}
	if !isValidSegment(sc.Segment()) {
		return fmt.Errorf("permission `%s`: invalid service segment", permission)
	}
	// Segment 2: Resource.
	if !sc.Scan() {
		return fmt.Errorf("permission `%s`: missing resource segment", permission)
	}
	if sc.Wildcard() {
		if sc.Scan() {
			return fmt.Errorf("permission `%s`: only final segment can be wildcard", permission)
		}
		return nil
	}
	if !isValidSegment(sc.Segment()) {
		return fmt.Errorf("permission `%s`: invalid resource segment", permission)
	}
	// Segment 3: Verb.
	if !sc.Scan() {
		return fmt.Errorf("permission `%s`: missing verb segment", permission)
	}
	if sc.Wildcard() {
		if sc.Scan() {
			return fmt.Errorf("permission `%s`: only final segment can be wildcard", permission)
		}
		return nil
	}
	if !isValidSegment(sc.Segment()) {
		return fmt.Errorf("permission `%s`: invalid verb segment", permission)
	}
	// Segment 4? Invalid!
	if sc.Scan() {
		return fmt.Errorf("permission `%s`: too many segments", permission)
	}
	return nil
}

// Valid checks whether the provided permission is valid.
// See Validate for what constitutes a valid permission.
func Valid(permission string) bool {
	return Validate(permission) == nil
}

func isValidSegment(segment string) bool {
	if len(segment) == 0 {
		return false
	}
	for i, r := range segment {
		switch i {
		case 0:
			if !unicode.IsLower(r) {
				return false
			}
		default:
			if !unicode.In(r, unicode.Lower, unicode.Digit) {
				return false
			}
		}
	}
	return true
}
