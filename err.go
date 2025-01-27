// Package builder provides a fluent SQL query builder with support for multiple SQL dialects.
package builder

import (
	"errors"
)

// Error variables for common SQL builder error conditions.
// These errors are returned by various methods in the package to indicate
// specific error conditions during query construction.
var (
	// ErrEmptyCondition is returned when attempting to build a condition without
	// providing the necessary field, operator, or values.
	ErrEmptyCondition = errors.New("empty condition")

	// ErrEmptySQLType is returned when attempting to build a query without
	// specifying the SQL operation type (SELECT, INSERT, etc.).
	ErrEmptySQLType = errors.New("empty sql type")

	// ErrListIsNotEmpty is returned when there are accumulated errors during
	// query construction. This typically indicates invalid SQL syntax or
	// incompatible operations.
	ErrListIsNotEmpty = errors.New("there are some errors in SQL, please check your query")
)
