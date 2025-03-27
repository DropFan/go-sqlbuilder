// Package builder provides a fluent SQL query builder with support for multiple SQL dialects.
package builder

import (
	"strconv"
	"strings"
)

// Dialector defines the interface for SQL dialect-specific operations.
// It provides methods for handling placeholder styles (e.g., ? for MySQL, $n for PostgreSQL)
// and identifier escaping conventions for different SQL databases.
//
// Each database type (MySQL, PostgreSQL, SQLite) has its own implementation
// of this interface to handle its specific SQL syntax requirements.
type Dialector interface {
	// Escape returns the escaped version of the provided identifiers
	// according to the dialect's escaping rules.
	Escape(s ...string) string

	// Placeholder returns the parameter placeholder for the given index.
	// For MySQL it returns "?", for PostgreSQL it returns "$n" where n is the index.
	Placeholder(index int) string

	// GetEscapeChar returns the character used for escaping identifiers
	// in the specific SQL dialect (e.g., backtick for MySQL, double quote for PostgreSQL).
	GetEscapeChar() string
}

var (
	// mysqlDialector is the default MySQL dialect implementation
	mysqlDialector MysqlDialector
	// postgresDialector is the PostgreSQL dialect implementation
	postgresDialector PostgresqlDialector
	// sqliteDialector is the SQLite dialect implementation
	sqliteDialector SQLiteDialector
)

// MysqlDialector implements the Dialector interface for MySQL database.
// It provides MySQL-specific SQL syntax handling, including backtick escaping
// for identifiers and ? for parameter placeholders.
type MysqlDialector struct {
	Dialector
}

// PostgresqlDialector implements the Dialector interface for PostgreSQL database.
// It provides PostgreSQL-specific SQL syntax handling, including double quote escaping
// for identifiers and $n for parameter placeholders.
type PostgresqlDialector struct {
	Dialector
}

// SQLiteDialector implements the Dialector interface for SQLite database.
// It provides SQLite-specific SQL syntax handling, including double quote escaping
// for identifiers and ? for parameter placeholders.
type SQLiteDialector struct {
	Dialector
}

// Escape wraps MySQL identifiers with backticks and handles multiple identifiers
// by joining them with "`, `".
func (MysqlDialector) Escape(s ...string) string {
	if len(s) == 1 {
		return "`" + s[0] + "`"
	}
	str := strings.Join(s, "`, `")
	return "`" + strings.Trim(str, "`") + "`"
}

// GetEscapeChar returns the backtick character used for escaping MySQL identifiers.
func (MysqlDialector) GetEscapeChar() string {
	return "`"
}

// Placeholder returns "?" as the parameter placeholder for MySQL queries.
func (MysqlDialector) Placeholder(index int) string {
	return "?"
}

// Escape wraps PostgreSQL identifiers with double quotes and handles multiple identifiers
// by joining them with '", "'.
func (p PostgresqlDialector) Escape(s ...string) string {
	if len(s) == 1 {
		return `"` + s[0] + `"`
	}
	str := strings.Join(s, `", "`)
	return `"` + strings.Trim(str, `"`) + `"`
}

// GetEscapeChar returns the double quote character used for escaping PostgreSQL identifiers.
func (p PostgresqlDialector) GetEscapeChar() string {
	return `"`
}

// Placeholder returns "$n" as the parameter placeholder for PostgreSQL queries,
// where n is the parameter index.
func (p PostgresqlDialector) Placeholder(index int) string {
	return "$" + strconv.Itoa(index)
}

// Escape wraps SQLite identifiers with double quotes and handles multiple identifiers
// by joining them with '", "'.
func (s SQLiteDialector) Escape(strs ...string) string {
	if len(strs) == 1 {
		return `"` + strs[0] + `"`
	}
	str := strings.Join(strs, `", "`)
	return `"` + strings.Trim(str, `"`) + `"`
}

// GetEscapeChar returns the double quote character used for escaping SQLite identifiers.
func (s SQLiteDialector) GetEscapeChar() string {
	return `"`
}

// Placeholder returns "?" as the parameter placeholder for SQLite queries.
func (s SQLiteDialector) Placeholder(index int) string {
	return "?"
}
