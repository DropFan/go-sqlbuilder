package builder

import (
	"strconv"
	"strings"
)

// Dialector eg: ? for mysql, $index for postgresql
type Dialector interface {
	Escape(s string) string
	Placeholder(index int) string
}

var (
	mysqlDialector MysqlDialector
	postgresDialector PostgresqlDialector
)

// MysqlDialector ...
type MysqlDialector struct {
	Dialector
}

// Escape ...
func (m MysqlDialector) Escape(s string) string {
	return "`" + strings.Trim(s, "`") + "`"
}

// Placeholder ...
func (m MysqlDialector) Placeholder(index int) string {
	return "?"
}

// PostgresqlDialector ...
type PostgresqlDialector struct {
	Dialector
}

// Escape ...
func (p PostgresqlDialector) Escape(s string) string {
	return "\"" + strings.Trim(s, "\"") + "\""
}

// Placeholder ...
func (p PostgresqlDialector) Placeholder(index int) string {
	return "$" + strconv.Itoa(index)
}
