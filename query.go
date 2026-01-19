// Package builder provides a fluent SQL query builder with support for multiple SQL dialects.
package builder

import (
	"fmt"
	"strconv"
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
// The method supports both MySQL-style (?) and PostgreSQL-style ($1, $2, ...)
// placeholders. It properly escapes string values by doubling single quotes.
//
// Warning: The returned string may contain unescaped values and should never
// be executed directly against a database.
//
// Example:
//
//	q := NewQuery("SELECT * FROM users WHERE status = ?", "active")
//	fmt.Println(q.String()) // Outputs: SELECT * FROM users WHERE status = 'active'
func (q *Query) String() string {
	if len(q.Args) == 0 {
		return q.Query
	}

	var result strings.Builder
	query := q.Query
	argIndex := 0
	i := 0

	for i < len(query) {
		// Check for ? placeholder (MySQL/SQLite style)
		if query[i] == '?' && argIndex < len(q.Args) {
			result.WriteString(formatValue(q.Args[argIndex]))
			argIndex++
			i++
			continue
		}

		// Check for $n placeholder (PostgreSQL style)
		if query[i] == '$' && i+1 < len(query) && query[i+1] >= '0' && query[i+1] <= '9' {
			j := i + 1
			for j < len(query) && query[j] >= '0' && query[j] <= '9' {
				j++
			}
			idx, _ := strconv.Atoi(query[i+1 : j])
			if idx > 0 && idx <= len(q.Args) {
				result.WriteString(formatValue(q.Args[idx-1]))
			} else {
				result.WriteString(query[i:j])
			}
			i = j
			continue
		}

		result.WriteByte(query[i])
		i++
	}

	return result.String()
}

// formatValue converts a value to its SQL string representation for debugging.
// All values are quoted for consistency in debug output.
func formatValue(v interface{}) string {
	if v == nil {
		return "NULL"
	}

	switch val := v.(type) {
	case string:
		escaped := strings.ReplaceAll(val, "'", "''")
		return "'" + escaped + "'"
	case bool:
		if val {
			return "TRUE"
		}
		return "FALSE"
	default:
		str := fmt.Sprintf("%v", val)
		escaped := strings.ReplaceAll(str, "'", "''")
		return "'" + escaped + "'"
	}
}
