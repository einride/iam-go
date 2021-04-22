package iampermission

// HasWildcard reports whether an IAM permission contains a wildcard '*' segment.
func HasWildcard(permission string) bool {
	var sc Scanner
	sc.Init(permission)
	for sc.Scan() {
		if sc.Wildcard() {
			return true
		}
	}
	return false
}
