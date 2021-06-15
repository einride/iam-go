package iammember

// New returns a new IAM member with the provided kind and value.
func New(kind string, value string) string {
	if kind == "" || value == "" {
		return ""
	}
	return kind + ":" + value
}
