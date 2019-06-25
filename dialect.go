package builder

import (
	"strconv"
	"strings"
)

// Dialector eg: ? for mysql, $index for postgresql
type Dialector interface {
	Escape(s ...string) string
	Placeholder(index int) string
	GetEscapeChar() string
}

var (
	mysqlDialector    MysqlDialector
	postgresDialector PostgresqlDialector
)

// MysqlDialector ...
type MysqlDialector struct {
	Dialector
}

// Escape ...
func (MysqlDialector) Escape(s ...string) string {
	str := strings.Join(s, "`, `")
	return "`" + strings.Trim(str, "`") + "`"
}

// GetEscapeChar ...
func (MysqlDialector) GetEscapeChar() string {
	return "`"
}

// Placeholder ...
func (MysqlDialector) Placeholder(index int) string {
	return "?"
}

// PostgresqlDialector ...
type PostgresqlDialector struct {
	Dialector
}

// Escape ...
func (p PostgresqlDialector) Escape(s ...string) string {
	str := strings.Join(s, `", "`)
	return `"` + strings.Trim(str, `"`) + `"`
}

// GetEscapeChar ...
func (p PostgresqlDialector) GetEscapeChar() string {
	return `"`
}

// Placeholder ...
func (p PostgresqlDialector) Placeholder(index int) string {
	return "$" + strconv.Itoa(index)
}
