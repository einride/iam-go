package iampermission

import "strings"

// Scanner scans an IAM permission.
type Scanner struct {
	permission string
	start, end int
}

// Init initializes the scanner.
func (s *Scanner) Init(permission string) {
	s.permission = permission
	s.start, s.end = 0, 0
}

// Scan to the next segment.
func (s *Scanner) Scan() bool {
	switch s.end {
	case len(s.permission):
		return false
	case 0:
		// start at the beginning
	default:
		s.start = s.end + 1 // start past latest dot '.'
	}
	if nextDot := strings.IndexByte(s.permission[s.start:], '.'); nextDot == -1 {
		s.end = len(s.permission)
	} else {
		s.end = s.start + nextDot
	}
	return true
}

// Segment returns the current segment.
func (s *Scanner) Segment() string {
	return s.permission[s.start:s.end]
}

// Wildcard reports whether the current segment is a wildcard.
func (s *Scanner) Wildcard() bool {
	return s.Segment() == "*"
}
