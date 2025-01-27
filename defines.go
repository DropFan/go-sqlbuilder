// Package builder provides a fluent SQL query builder with support for multiple SQL dialects.
package builder

// SQLType represents the type of SQL query being constructed.
// It is used to track and validate the query type during construction.
type SQLType int

// SQL query types supported by the builder.
const (
	// RawSQL represents a raw SQL query that will be used as-is.
	RawSQL SQLType = iota + 1

	// SelectSQL represents a SELECT query for retrieving data.
	SelectSQL

	// InsertSQL represents an INSERT query for adding new records.
	InsertSQL

	// ReplaceSQL represents a REPLACE query (MySQL-specific) for replacing records.
	ReplaceSQL

	// InsertOrUpdateSQL represents an INSERT ... ON DUPLICATE KEY UPDATE query (MySQL-specific).
	InsertOrUpdateSQL

	// UpdateSQL represents an UPDATE query for modifying existing records.
	UpdateSQL

	// DeleteSQL represents a DELETE query for removing records.
	DeleteSQL
)

// operMap defines the mapping between SQL operators and their expected number of values.
// A value of 1 indicates a single-value operator (e.g., =, >).
// A value of 2 indicates a two-value operator (e.g., BETWEEN).
// A value of 3 indicates a multi-value operator (e.g., IN).
var operMap = map[string]int{
	"=":           1, // Equal
	"!=":          1, // Not equal
	"<>":          1, // Not equal (alternative syntax)
	">":           1, // Greater than
	">=":          1, // Greater than or equal
	"<":           1, // Less than
	"<=":          1, // Less than or equal
	"IN":          3, // IN clause (multiple values)
	"NOT IN":      3, // NOT IN clause (multiple values)
	"LIKE":        1, // Pattern matching
	"NOT LIKE":    1, // Negative pattern matching
	"BETWEEN":     2, // Range comparison
	"NOT BETWEEN": 2, // Negative range comparison
}
