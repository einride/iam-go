package iamexample

import _ "embed"

//go:embed schema.sql
var spannerSQLSchema string

// SQLSchema returns the example Spanner SQL schema.
func SQLSchema() string {
	return spannerSQLSchema
}
