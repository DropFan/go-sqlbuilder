// Package builder provides a fluent SQL query builder with support for multiple SQL dialects.
// It offers a clean and intuitive API for constructing SQL queries programmatically while
// handling proper escaping and parameter binding based on the target database system.
//
// The package supports multiple SQL dialects including MySQL, PostgreSQL, and SQLite,
// with appropriate escaping and placeholder styles for each. It provides a type-safe
// way to build complex SQL queries without string concatenation or manual escaping.
//
// Example usage:
//
//	b := builder.New()
//	query, err := b.Select("id", "name").From("users").Where(builder.Eq("status", "active")).Build()
//
// For more examples, see the builder_test.go file.
package builder

import (
	"fmt"
	"strconv"
	"strings"
)

// Builder represents a SQL query builder that supports multiple SQL dialects.
// It provides a fluent interface for constructing SQL queries with proper escaping
// and parameter binding.
//
// It is recommended to create a new Builder instance using the New() function.
// For usage examples, please refer to builder_test.go.
type Builder struct {
	// sqlType represents the type of SQL query being built (SELECT, INSERT, etc.)
	sqlType SQLType
	// dialector handles database-specific SQL syntax and escaping rules
	dialector Dialector
	// queryArgs holds the arguments for parameter binding in the query
	queryArgs []interface{}
	// query stores the current SQL query string being constructed using strings.Builder for better performance
	query strings.Builder
	// setValues stores the field names being updated in an UPDATE query
	setValues []string
	// ErrList collects any errors encountered during query construction
	ErrList []error
	// lastQueries maintains a history of all queries built by this instance
	lastQueries []*Query
	// The following fields are deprecated and will be removed in a future version:

	// queryTables string // abandoned
	// from    string     // abandoned
	// where []string     // abandoned
	// orderBy string     // abandoned
	// limit string       // abandoned
}

// SetDialector sets the SQL dialect for parameter binding placeholders and identifier escaping.
// It returns the Builder instance for method chaining.
func (b *Builder) SetDialector(d Dialector) *Builder {
	b.dialector = d
	return b
}

// Escape escapes the provided field names according to the current SQL dialect's rules.
// It returns the escaped string representation of the fields.
func (b *Builder) Escape(s ...string) string {
	return b.dialector.Escape(s...)
}

// EscapeChar returns the escape character used by the current SQL dialect
// for escaping identifiers.
func (b *Builder) EscapeChar() string {
	return b.dialector.GetEscapeChar()
}

// Placeholder returns the appropriate parameter placeholder for the current SQL dialect
// at the given position. For example, MySQL uses "?" while PostgreSQL uses "$1", "$2", etc.
//
// This method is currently disabled as the placeholder handling is done internally.
// func (b *Builder) Placeholder(index int) string {
// 	return b.dialector.Placeholder(index)
// }

// New creates and initializes a new Builder instance with default MySQL dialect.
// For usage examples, please refer to builder_test.go.
func New() *Builder {
	return &Builder{
		sqlType:   0,
		dialector: mysqlDialector,
		// queryTables: "",
		queryArgs: []interface{}{},
		query:     strings.Builder{},
		// from:      "",
		setValues: []string{},
		// where: []string{},
		// orderBy:     "",
		// limit:       "",
		ErrList:     []error{},
		lastQueries: []*Query{},
	}
}

// LastQueries returns all previously built queries in this builder instance.
func (b *Builder) LastQueries() []*Query {
	return b.lastQueries
}

// LastQuery returns the most recently built query, or nil if no queries have been built.
func (b *Builder) LastQuery() *Query {
	if len(b.lastQueries) == 0 {
		return nil
	}

	return b.lastQueries[len(b.lastQueries)-1]
}

// renew reset some data after `Build()` was called
func (b *Builder) renew(st SQLType) {
	if len(b.ErrList) > 0 {
		b.ErrList = b.ErrList[:0]
	} else {
		b.ErrList = []error{}
	}
	b.sqlType = st
	if len(b.queryArgs) > 0 {
		b.queryArgs = b.queryArgs[:0]
	} else {
		b.queryArgs = []interface{}{}
	}
	b.query.Reset()
	if len(b.setValues) > 0 {
		b.setValues = b.setValues[:0]
	} else {
		b.setValues = []string{}
	}
	// if len(b.where) > 0 {
	// 	b.where = b.where[:0]
	// } else {
	// 	b.where = []string{}
	// }
	// b.orderBy = ""
	// b.limit = ""
}

// Clear resets the current query and its arguments to their initial state.
// It returns the Builder instance for method chaining.
func (b *Builder) Clear() *Builder {
	b.renew(0)
	return b
}

// QueryArgs returns the current list of query arguments that will be used
// for parameter binding.
func (b *Builder) QueryArgs() []interface{} {
	return b.queryArgs
}

// Query returns the current SQL query string being constructed.
func (b *Builder) Query() string {
	return b.query.String()
}

// Append adds the provided string and arguments to the end of the current query.
// It returns the Builder instance for method chaining.
func (b *Builder) Append(s string, args ...interface{}) *Builder {
	b.query.WriteString(s)
	b.queryArgs = append(b.queryArgs, args...)
	return b
}

// AppendPre adds the provided string and arguments to the beginning of the current query.
// It returns the Builder instance for method chaining.
func (b *Builder) AppendPre(s string, args ...interface{}) *Builder {
	oldQuery := b.query.String()
	b.query.Reset()
	b.query.WriteString(s)
	b.query.WriteString(oldQuery)
	b.queryArgs = append(args, b.queryArgs...)
	return b
}

// Raw sets a raw SQL query string with optional arguments.
// It returns the Builder instance for method chaining.
func (b *Builder) Raw(s string, args ...interface{}) *Builder {
	b.renew(RawSQL)
	b.query.WriteString(s)
	b.queryArgs = append(args, b.queryArgs...)
	return b
}

// Select begins a SELECT query with the specified fields.
// If no fields are provided, it creates an empty SELECT.
// If "*" is provided as the first field, it selects all columns.
// It returns the Builder instance for method chaining.
func (b *Builder) Select(fields ...string) *Builder {
	b.renew(SelectSQL)
	b.query.WriteString("SELECT")

	if len(fields) <= 0 {
		// Do nothing
	} else if fields[0] == "*" {
		b.query.WriteString(" *")
	} else {
		b.query.WriteString(" ")
		b.query.WriteString(b.Escape(fields...))
		// b.query += " `" + strings.Join(fields, "`, `") + "`"
	}

	return b
}

// Insert begins an INSERT query for the specified table and optional field names.
// It returns the Builder instance for method chaining.
func (b *Builder) Insert(tableName string, fields ...string) *Builder {
	b.renew(InsertSQL)
	b.query.WriteString("INSERT INTO ")
	b.query.WriteString(b.Escape(tableName))

	if len(fields) > 0 {
		b.Into(fields...)
	}
	return b
}

// InsertOrUpdate begins an INSERT ... ON DUPLICATE KEY UPDATE query for MySQL.
// It takes field-value pairs that will be used for both the INSERT and UPDATE parts.
// It returns the Builder instance for method chaining.
func (b *Builder) InsertOrUpdate(tableName string, fvals ...*FieldValue) *Builder {
	b.renew(InsertSQL)
	b.query.WriteString("INSERT INTO ")
	b.query.WriteString(b.Escape(tableName))

	if len(fvals) > 0 {
		var (
			fields []string
			vals   []interface{}
		)
		for _, fv := range fvals {
			if fv != nil {
				fields = append(fields, fv.Name)
				vals = append(vals, fv.Value)
			}
		}
		b.Into(fields...).Values(vals).Append(" ON DUPLICATE KEY UPDATE ").Set(fvals...)
	}

	return b
}

// Replace begins a REPLACE query for the specified table and optional field names.
// It returns the Builder instance for method chaining.
func (b *Builder) Replace(tableName string, fields ...string) *Builder {
	b.renew(InsertSQL)
	b.query.WriteString("REPLACE INTO ")
	b.query.WriteString(b.Escape(tableName))

	if len(fields) > 0 {
		b.Into(fields...)
	}
	return b
}

// Into specifies the fields for an INSERT or REPLACE query.
// It returns the Builder instance for method chaining.
func (b *Builder) Into(fields ...string) *Builder {
	b.query.WriteString(" (")
	b.query.WriteString(b.Escape(fields...))
	b.query.WriteString(")")
	// b.query += " (`" + strings.Join(fields, "`, `") + "`)"
	return b
}

// Predefine placeholder strings for common quantities to avoid runtime calculations.
// Use the predefined placeholder string when there are less than 6 values.
var __placeholders = []string{
	"(?)",
	"(?, ?)",
	"(?, ?, ?)",
	"(?, ?, ?, ?)",
	"(?, ?, ?, ?, ?)",
}

// Values adds one or more sets of values to an INSERT or REPLACE query.
// Each set of values must match the number of fields specified in Into().
// It returns the Builder instance for method chaining.
func (b *Builder) Values(valsGroup ...[]interface{}) *Builder {
	b.query.WriteString(" VALUES ")
	// index := 0
	for i, vals := range valsGroup {
		if i > 0 {
			b.query.WriteString(", ")
		}
		// b.query += "("
		// for j, val := range vals {
		// 	index++
		// 	if j > 0 {
		// 		b.query += ", "
		// 	}
		// 	b.query += b.Placeholder(index)
		// 	b.queryArgs = append(b.queryArgs, val)
		// }
		// b.query += ")"

		// Use the predefined placeholder string when there are less than 6 values.
		if len(vals) > 5 {
			b.query.WriteString("(?")
			b.query.WriteString(strings.Repeat(", ?", len(vals)-1))
			b.query.WriteString(")")
		} else {
			b.query.WriteString(__placeholders[len(vals)-1])
		}
		b.queryArgs = append(b.queryArgs, vals...)
	}
	return b
}

// Update begins an UPDATE query for the specified table with optional field-value pairs.
// It returns the Builder instance for method chaining.
func (b *Builder) Update(tableName string, fvals ...*FieldValue) *Builder {
	b.renew(UpdateSQL)
	b.query.WriteString("UPDATE ")
	b.query.WriteString(b.Escape(tableName))
	b.query.WriteString(" SET ")

	if len(fvals) > 0 {
		b.Set(fvals...)
	}

	return b
}

// Set specifies the field-value pairs to update in an UPDATE query.
// It returns the Builder instance for method chaining.
func (b *Builder) Set(fvals ...*FieldValue) *Builder {
	// b.setValue = ""

	for i, fval := range fvals {
		if fval == nil {
			continue
		}
		if i > 0 || len(b.setValues) > 0 {
			b.query.WriteString(", ")
		}
		b.setValues = append(b.setValues, fval.Name)
		b.query.WriteString(b.Escape(fval.Name))
		b.query.WriteString(" = ?")
		b.queryArgs = append(b.queryArgs, fval.Value)
	}

	return b
}

// Delete begins a DELETE query for the specified table.
// It returns the Builder instance for method chaining.
func (b *Builder) Delete(tableName string) *Builder {
	b.renew(DeleteSQL)
	b.query.WriteString("DELETE FROM ")
	b.query.WriteString(b.Escape(tableName))

	return b
}

// Build finalizes the query construction and returns a Query object along with any errors.
// It validates the SQL type and any accumulated errors before creating the final query.
func (b *Builder) Build(queries ...interface{}) (q *Query, err error) {

	switch b.sqlType {
	case SelectSQL:
	case InsertSQL:
	// case InsertOrUpdateSQL:
	case UpdateSQL:
	case DeleteSQL:
	case RawSQL:
	default:
		return nil, ErrEmptySQLType
	}
	if len(b.ErrList) > 0 {
		err = ErrListIsNotEmpty
	}
	q = NewQuery(b.query.String(), b.queryArgs...)
	b.lastQueries = append(b.lastQueries, q)
	b.renew(RawSQL)
	return q, err
}

// From specifies the tables to select from in a SELECT query.
// It returns the Builder instance for method chaining.
func (b *Builder) From(tables ...string) *Builder {
	if len(tables) <= 0 {
		return b
	}
	// b.Tables = tables
	// b.QueryTables = "`" + strings.Join(tables, "`, `") + "`"
	b.query.WriteString(" FROM ")
	b.query.WriteString(b.Escape(tables...))
	// b.query += " FROM `" + strings.Join(tables, "`, `") + "`"
	return b
}

// FromRaw specifies a raw FROM clause without any escaping.
// It returns the Builder instance for method chaining.
func (b *Builder) FromRaw(from string) *Builder {
	b.query.WriteString(" FROM ")
	b.query.WriteString(from)
	return b
}

// Where begins the WHERE clause of a query with the specified conditions.
// If no conditions are provided, it adds "WHERE 1".
// It returns the Builder instance for method chaining.
// addConditions adds one or more conditions to the query.
// It handles the AND/OR logic between conditions and builds each condition
// using buildCondition. Any errors encountered during condition building
// are collected in the Builder's ErrList.
//
// Parameters:
//   - conditions: Variable number of Condition objects to be added
//
// Returns:
//   - *Builder: The Builder instance for method chaining
//
// Example:
//
//	b.addConditions(
//	  builder.Eq("status", "active"),
//	  builder.Gt("age", 18)
//	)
//	// Generates: `status` = ? AND `age` > ?
func (b *Builder) addConditions(conditions ...*Condition) *Builder {
	condSlice := make([]string, 0, len(conditions))
	for i, cond := range conditions {
		if cond == nil {
			continue
		}
		condStr, args, err := b.buildCondition(cond)
		if err != nil {
			condStr = fmt.Sprintf("{error: %s}", err)
			b.ErrList = append(b.ErrList, err)
		}
		if cond.AndOr && i > 0 {
			condStr = "AND " + condStr
		} else if i > 0 {
			condStr = "OR " + condStr
		}
		condSlice = append(condSlice, condStr)
		b.queryArgs = append(b.queryArgs, args...)
	}
	// b.where = append(b.where, "("+strings.Join(condSlice, " ")+")")
	// if len(conditions) <= 0 {
	// 	condSlice = append(condSlice, "1")
	// }
	// b.where = append(b.where, strings.Join(condSlice, " "))
	b.query.WriteString(strings.Join(condSlice, " "))
	return b

}

// And adds one or more conditions to the query using AND logic.
// If a single condition is provided, it adds "AND condition".
// If multiple conditions are provided, it adds "AND (condition1 AND condition2 ...)"
// It returns the Builder instance for method chaining.
func (b *Builder) And(conditions ...*Condition) *Builder {
	switch len(conditions) {
	case 0:
		return b
	case 1:
		b.query.WriteString(" AND ")
		b.addConditions(conditions...)
	default:
		b.query.WriteString(" AND (")
		b.addConditions(conditions...)
		b.query.WriteString(")")
	}
	// if len(b.where) > 0 {
	// 	b.where[0] = "(" + b.where[0]
	// 	b.where[len(b.where)-1] = b.where[len(b.where)-1] + ") AND"
	// }
	// b.where = append(b.where, ") AND")
	// b.where = append(b.where, ")")
	return b
}

// Or adds one or more conditions to the query using OR logic.
// If a single condition is provided, it adds "OR condition".
// If multiple conditions are provided, it adds "OR (condition1 OR condition2 ...)"
// It returns the Builder instance for method chaining.
func (b *Builder) Or(conditions ...*Condition) *Builder {
	switch len(conditions) {
	case 0:
		return b
	case 1:
		b.query.WriteString(" OR ")
		b.addConditions(conditions...)
	default:
		b.query.WriteString(" OR (")
		b.addConditions(conditions...)
		b.query.WriteString(")")
	}

	// if len(b.where) > 0 {
	// 	b.where[0] = "(" + b.where[0]
	// 	b.where[len(b.where)-1] = b.where[len(b.where)-1] + ") OR"
	// }
	// b.addConditions(conditions...)
	// b.where = append(b.where, "")
	return b
}

// In adds an IN condition to the query for the specified field and values.
// It is equivalent to "field IN (value1, value2, ...)"
// It returns the Builder instance for method chaining.
//
// Example:
//
//	b.Select("*").From("users").In("status", "active", "pending")
//	// Generates: SELECT * FROM users WHERE status IN (?, ?)
//
// In creates an IN condition for the specified field and values.
// It returns the Builder instance for method chaining.
//
// Parameters:
//   - field: The field name to check against
//   - values: Variable number of values to include in the IN clause
//
// Example:
//
//	b.Select("*").From("users")
//	  .Where(builder.In("status", "active", "pending"))
//	// Generates: SELECT * FROM users WHERE `status` IN (?, ?)
func (b *Builder) In(field string, values ...interface{}) *Builder {
	b.addConditions(In(field, values...))
	return b
}

// NotIn adds a NOT IN condition to the query for the specified field and values.
// It is equivalent to "field NOT IN (value1, value2, ...)"
// It returns the Builder instance for method chaining.
//
// Example:
//
//	b.Select("*").From("users").NotIn("status", "deleted", "banned")
//	// Generates: SELECT * FROM users WHERE status NOT IN (?, ?)
func (b *Builder) NotIn(field string, values ...interface{}) *Builder {
	b.addConditions(NotIn(field, values...))
	return b
}

// Between adds a BETWEEN condition to the query for the specified field and range values.
// It expects exactly two values defining the range (start and end).
// It returns the Builder instance for method chaining.
//
// Example:
//
//	b.Select("*").From("orders").Between("amount", 100, 1000)
//	// Generates: SELECT * FROM orders WHERE amount BETWEEN ? AND ?
func (b *Builder) Between(field string, values ...interface{}) *Builder {
	b.addConditions(Between(field, values...))
	return b
}

// NotBetween adds a NOT BETWEEN condition to the query for the specified field and range values.
// It expects exactly two values defining the range (start and end).
// It returns the Builder instance for method chaining.
//
// Example:
//
//	b.Select("*").From("orders").NotBetween("amount", 0, 100)
//	// Generates: SELECT * FROM orders WHERE amount NOT BETWEEN ? AND ?
func (b *Builder) NotBetween(field string, values ...interface{}) *Builder {
	b.addConditions(NotBetween(field, values...))
	return b
}

// Where begins the WHERE clause of a query with the specified conditions.
// If no conditions are provided, it adds "WHERE 1".
// It returns the Builder instance for method chaining.
func (b *Builder) Where(conditions ...*Condition) *Builder {
	b.query.WriteString(" WHERE ")
	if len(conditions) == 0 {
		b.query.WriteString("1")
		return b
	}

	b.addConditions(conditions...)
	return b
}

// WhereRaw adds a raw WHERE clause without any escaping or parameter binding.
// This method is useful when you need to write complex WHERE conditions that
// cannot be easily expressed using the standard condition builders.
// It returns the Builder instance for method chaining.
//
// Warning: Be careful when using this method with user-provided input as it
// may lead to SQL injection vulnerabilities. Use the standard Where method
// with proper parameter binding when possible.
//
// Parameters:
//   - str: The raw WHERE clause string
//   - args: Optional arguments for parameter binding
//
// Example:
//
//	// Using raw SQL function
//	b.Select("*").From("posts")
//	  .WhereRaw("DATE(created_at) = CURDATE()")
//
//	// With parameter binding
//	b.Select("*").From("users")
//	  .WhereRaw("FIND_IN_SET(?, roles)", "admin")
//	// Generates: SELECT * FROM users WHERE FIND_IN_SET(?, roles)
func (b *Builder) WhereRaw(str string, args ...interface{}) *Builder {
	b.query.WriteString(" WHERE ")
	b.query.WriteString(str)
	b.queryArgs = append(b.queryArgs, args...)

	return b
}

// buildCondition constructs a SQL condition string and its corresponding query arguments.
// It handles various SQL operators (=, !=, IN, BETWEEN, etc.) and validates their usage.
// Returns the condition string, query arguments, and any validation errors.
//
// The method validates:
// - Operator existence in the operMap
// - Number of values matching the operator requirements
// - Proper formatting of placeholders based on operator type
// buildCondition constructs a SQL condition string and its corresponding query arguments.
// It handles various SQL operators (=, !=, IN, BETWEEN, etc.) and validates their usage.
// Returns the condition string, query arguments, and any validation errors.
//
// Parameters:
//   - cond: A Condition object containing the field, operator, and values to build the condition
//
// Returns:
//   - str: The constructed SQL condition string
//   - queryArgs: Slice of arguments for parameter binding
//   - err: Error if validation fails
//
// Supported operators:
//   - Single value: =, !=, <>, >, >=, <, <=, LIKE, NOT LIKE
//   - Multiple values: IN, NOT IN
//   - Range values: BETWEEN, NOT BETWEEN
//
// Example:
//
//	cond := &Condition{Field: "age", Operator: ">", Values: []interface{}{18}}
//	str, args, err := builder.buildCondition(cond)
//	// Generates: `age` > ? with args=[18]
func (b *Builder) buildCondition(cond *Condition) (str string, queryArgs []interface{}, err error) {
	// if cond == nil {
	// 	return "", nil, ErrEmptyCondition
	// }

	str = ""
	queryArgs = []interface{}{}

	if opValue, ok := operMap[cond.Operator]; !ok {
		// return "", queryArgs,
		err = fmt.Errorf("invalid operator:(operator:%s[field:%s])", cond.Operator, cond.Field)
		return
	} else if len(cond.Values) != opValue {
		switch opValue {
		case 3:
			if len(cond.Values) != 0 {
				break
			}
			fallthrough
		case 1, 2:
			err = fmt.Errorf("invalid number of values with operator:(%s[field:%s])", cond.Operator, cond.Field)
			return
		}
	}
	// "=":           1,
	// "!=":          1,
	// "<>":          1,
	// ">":           1,
	// ">=":          1,
	// "<":           1,
	// "<=":          1,
	// "IN":          3,
	// "NOT IN":      3,
	// "LIKE":        1,
	// "NOT LIKE":    1,
	// "BETWEEN":     2,
	// "NOT BETWEEN": 2,
	placeholders := ""
	switch strings.ToLower((cond.Operator)) {
	case "=",
		"!=", "<>",
		">", ">=",
		"<", "<=",
		"like", "not like":
		placeholders = "?"
	case "in", "not in":
		placeholders = "(?" + strings.Repeat(", ?", len(cond.Values)-1) + ")"
	case "between", "not between":
		placeholders += "? AND ?"
		// default:
	}

	queryArgs = append(queryArgs, cond.Values...)

	str += b.Escape(cond.Field) + " " + cond.Operator + " " + placeholders

	return
}

// OrderBy specifies the ORDER BY clause with the given conditions.
// Each condition determines the field and sort direction (ASC/DESC).
// Multiple conditions can be combined to sort by multiple fields.
// It returns the Builder instance for method chaining.
//
// Parameters:
//   - conditions: One or more Condition objects specifying the sort fields and directions
//
// Example:
//
//	b.Select("*").From("users")
//	  .OrderBy(
//	    builder.Asc("created_at"),    // Sort by created_at ASC
//	    builder.Desc("last_login"),   // Then by last_login DESC
//	  )
//	// Generates: SELECT * FROM users ORDER BY `created_at` ASC, `last_login` DESC
func (b *Builder) OrderBy(conditions ...*Condition) *Builder {
	b.query.WriteString(" ORDER BY ")

	condStrSlice := []string{}

	for _, cond := range conditions {
		if cond == nil {
			continue
		}
		var condStr strings.Builder
		condStr.WriteString(b.Escape(cond.Field))
		if cond.Asc {
			condStr.WriteString(" ASC")
		} else {
			condStr.WriteString(" DESC")
		}
		condStrSlice = append(condStrSlice, condStr.String())
	}
	b.query.WriteString(strings.Join(condStrSlice, ", "))

	return b
}

// Limit adds a LIMIT clause to the query to restrict the number of rows returned.
// It can be used in two ways:
// 1. With a single argument to limit the number of rows
// 2. With two arguments to specify both offset and limit
//
// Parameters:
//   - limitOffset: One or two integers:
//   - With one argument: The maximum number of rows to return
//   - With two arguments: The offset and the maximum number of rows
//
// Example:
//
//	// Limit to 10 rows
//	b.Select("*").From("users").Limit(10)
//	// Skip 20 rows and return next 10
//	b.Select("*").From("users").Limit(20, 10)
func (b *Builder) Limit(limitOffset ...int) *Builder {
	if len(limitOffset) == 1 {
		b.query.WriteString(" LIMIT ")
		b.query.WriteString(strconv.Itoa(limitOffset[0]))
	} else {
		b.query.WriteString(" LIMIT ")
		b.query.WriteString(strconv.Itoa(limitOffset[0]))
		b.query.WriteString(", ")
		b.query.WriteString(strconv.Itoa(limitOffset[1]))
	}
	return b
}

// Count begins a SELECT COUNT query.
// It can count all rows (COUNT(1)) or specific fields/expressions.
// It returns the Builder instance for method chaining.
//
// Parameters:
//   - query: Optional string(s) to specify what to count (e.g., "DISTINCT column")
//     If not provided, defaults to COUNT(1)
//
// Example:
//
//	// Count all rows
//	b.Count().From("users")
//	// Generates: SELECT COUNT(1) FROM users
//
//	// Count distinct values
//	b.Count("DISTINCT status").From("orders")
//	// Generates: SELECT COUNT(DISTINCT status) FROM orders
func (b *Builder) Count(query ...string) *Builder {
	b.renew(SelectSQL)
	b.query.WriteString("SELECT COUNT(")
	if len(query) <= 0 {
		b.query.WriteString("1")
	} else {
		b.query.WriteString(strings.TrimSpace(strings.Join(query, " ")))
	}
	b.query.WriteString(")")

	return b
}
