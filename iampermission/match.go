package iampermission

// Match reports whether the lhs pattern name matches the specified resource name pattern.
// The result will always be false when any of lhs or rhs are invalid, or if rhs contains a wildcard.
func Match(lhs, rhs string) bool {
	if !Valid(lhs) || !Valid(rhs) || HasWildcard(rhs) {
		return false
	}
	var scLhs Scanner
	scLhs.Init(lhs)
	var scRhs Scanner
	scRhs.Init(rhs)
	// Segment 1: Service.
	_, _ = scLhs.Scan(), scRhs.Scan()
	if scLhs.Segment() != scRhs.Segment() {
		return false
	}
	// Segment 2: Resource.
	_, _ = scLhs.Scan(), scRhs.Scan()
	if scLhs.Wildcard() {
		return true
	}
	if scLhs.Segment() != scRhs.Segment() {
		return false
	}
	// Segment 3: Verb.
	_, _ = scLhs.Scan(), scRhs.Scan()
	if scLhs.Wildcard() {
		return true
	}
	return scLhs.Segment() == scRhs.Segment()
}
