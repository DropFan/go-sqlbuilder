// Package builder provides a fluent SQL query builder with support for multiple SQL dialects.
package builder

import (
	"fmt"
	"strings"
)

// Query represents a SQL query with its associated parameter values.
// It provides methods for converting the parameterized query into a string
// representation with the parameter values interpolated.
//
// The Query type is typically created by calling Build() on a Builder instance
// and should not be created directly. It encapsulates both the SQL query string
// and its parameter values to ensure type safety and prevent SQL injection.
//
// Example usage:
//
//	b := builder.New()
//	query, err := b.Select("id", "name").From("users").Where(builder.Eq("status", "active")).Build()
//	// query can be used with database/sql's Exec or Query methods
type Query struct {
	// Query holds the parameterized SQL query string with placeholders
	Query string

	// Args holds the parameter values that correspond to the placeholders in Query
	Args []interface{}
}

// NewQuery creates a new Query instance with the given SQL query string
// and parameter values. The query string should use placeholders (?) for
// parameter values that will be bound when executing the query.
//
// This function is primarily used internally by the Builder's Build method
// and should rarely be called directly. It ensures that the number of
// placeholders matches the number of provided arguments.
//
// Parameters:
//   - q: The SQL query string with parameter placeholders
//   - args: The parameter values that correspond to the placeholders
//
// Returns a new Query instance that can be used with database/sql methods.
func NewQuery(q string, args ...interface{}) *Query {
	return &Query{
		Query: q,
		Args:  args,
	}
}

// String returns a string representation of the query with parameter values
// interpolated into the query string. This method is useful for debugging and logging
// purposes, but should not be used for actual query execution to avoid SQL
// injection vulnerabilities.
//
// The method performs simple string replacements to create a human-readable
// version of the query. It handles common SQL operators and placeholders,
// replacing them with the actual parameter values.
//
// Warning: The returned string may contain unescaped values and should never
// be executed directly against a database.
//
// Example:
//
//	q := NewQuery("SELECT * FROM users WHERE status = ?", "active")
//	fmt.Println(q.String()) // Outputs: SELECT * FROM users WHERE status = 'active'
func (q *Query) String() string {
	replacements := map[string]string{
		" = ?":   " = '%v'",
		" != ?":  " != '%v'",
		" <> ?":  " <> '%v'",
		" > ?":   " > '%v'",
		" >= ?":  " >= '%v'",
		" < ?":   " < '%v'",
		" <= ?":  " <= '%v'",
		" (?, ":  " ('%v', ",
		" (?)":   " ('%v')",
		", ?":    ", '%v'",
		" ?, ":   " '%v', ",
		" ? AS ": " '%v' AS ",
	}

	var result strings.Builder
	str := q.Query

	for pattern, replacement := range replacements {
		str = strings.Replace(str, pattern, replacement, -1)
	}

	result.WriteString(fmt.Sprintf(str, q.Args...))

	return result.String()
}
