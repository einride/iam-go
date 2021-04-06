package iampermission

// Match reports whether the lhs permission name matches the rhs permission.
// The lhs permission may contain a wildcard.
// The result will always be false when any of lhs or rhs are invalid, or if rhs contains a wildcard.
func Match(lhs, rhs string) bool {
	if !Valid(lhs) || !Valid(rhs) || HasWildcard(rhs) {
		return false
	}
	var scLHS Scanner
	scLHS.Init(lhs)
	var scRHS Scanner
	scRHS.Init(rhs)
	// Segment 1: Service.
	_, _ = scLHS.Scan(), scRHS.Scan()
	if scLHS.Segment() != scRHS.Segment() {
		return false
	}
	// Segment 2: Resource.
	_, _ = scLHS.Scan(), scRHS.Scan()
	if scLHS.Wildcard() {
		return true
	}
	if scLHS.Segment() != scRHS.Segment() {
		return false
	}
	// Segment 3: Verb.
	_, _ = scLHS.Scan(), scRHS.Scan()
	if scLHS.Wildcard() {
		return true
	}
	return scLHS.Segment() == scRHS.Segment()
}
