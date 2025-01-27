// Package builder provides a fluent SQL query builder with support for multiple SQL dialects.
package builder

// Values represents a slice of interface{} used for storing query parameter values.
// It is used throughout the package to handle various types of SQL parameter values
// in a type-safe manner. This type provides a convenient way to pass multiple
// values to SQL conditions while maintaining type flexibility.
//
// Example usage:
//
//	values := Values{1, "active", true}
//	// Can be used with IN conditions
//	b.Select("*").From("users").Where(builder.In("status", values...))
type Values []interface{}

// Condition represents a SQL condition that can be used in WHERE clauses or ORDER BY statements.
// It supports various SQL operators and can be combined using AND/OR logic.
//
// Conditions can be created using helper functions like And(), Or(), Between(), In(), etc.
// Multiple conditions can be combined to create complex WHERE clauses with proper
// operator precedence and grouping.
//
// Example usage:
//
//	// Simple condition
//	cond1 := builder.And("status", "=", "active")
//
//	// Multiple conditions combined
//	cond2 := builder.Or("role", "IN", "admin", "moderator")
//	b.Select("*").From("users").Where(cond1, cond2)
type Condition struct {
	Field    string        // The name of the database field or column
	Asc      bool          // Sort direction: true for ASC, false for DESC
	AndOr    bool          // Logical operator: true for AND, false for OR
	Operator string        // SQL operator (e.g., =, >, LIKE, IN, etc.)
	Values   []interface{} // The values to compare against the field
}

// newCondition creates a new Condition with the specified parameters.
// It is an internal helper function used by the public condition constructors.
func newCondition(andOr bool, field string, op string, values []interface{}) *Condition {
	return &Condition{
		AndOr:    andOr,
		Field:    field,
		Operator: op,
		Values:   values,
	}
}

// And creates a new condition that will be combined with AND logic.
// It takes a field name, operator, and values to construct the condition.
//
// Parameters:
//   - field: The database column or field name
//   - op: The SQL operator (e.g., "=", ">", "LIKE", etc.)
//   - values: The values to compare against the field
//
// Example:
//
//	// Creates: WHERE age >= 18 AND status = 'active'
//	b.Where(
//		builder.And("age", ">=", 18),
//		builder.And("status", "=", "active"),
//	)
func And(field string, op string, values ...interface{}) *Condition {
	return newCondition(true, field, op, values)
}

// Or creates a new condition that will be combined with OR logic.
// It takes a field name, operator, and values to construct the condition.
//
// Parameters:
//   - field: The database column or field name
//   - op: The SQL operator (e.g., "=", ">", "LIKE", etc.)
//   - values: The values to compare against the field
//
// Example:
//
//	// Creates: WHERE role = 'admin' OR role = 'moderator'
//	b.Where(
//		builder.Or("role", "=", "admin"),
//		builder.Or("role", "=", "moderator"),
//	)
func Or(field string, op string, values ...interface{}) *Condition {
	return newCondition(false, field, op, values)
}

// Between creates a new BETWEEN condition for the specified field.
// The values parameter should contain exactly two values defining the range.
//
// Parameters:
//   - field: The database column or field name
//   - values: Exactly two values defining the range (start and end)
//
// Example:
//
//	// Creates: WHERE price BETWEEN 10 AND 100
//	b.Where(builder.Between("price", 10, 100))
func Between(field string, values ...interface{}) *Condition {
	return newCondition(false, field, "BETWEEN", values)
}

// NotBetween creates a new NOT BETWEEN condition for the specified field.
// The values parameter should contain exactly two values defining the range.
//
// Parameters:
//   - field: The database column or field name
//   - values: Exactly two values defining the range (start and end)
//
// Example:
//
//	// Creates: WHERE age NOT BETWEEN 0 AND 18
//	b.Where(builder.NotBetween("age", 0, 18))
func NotBetween(field string, values ...interface{}) *Condition {
	return newCondition(false, field, "NOT BETWEEN", values)
}

// Eq creates a new equality condition for the specified field and value.
//
// Parameters:
//   - field: The database column or field name
//   - value: The value to compare against the field
//
// Example:
//
//	// Creates: WHERE status = 'active'
//	b.Where(builder.Eq("status", "active"))
func Eq(field string, value interface{}) *Condition {
	return newCondition(false, field, "=", []interface{}{value})
}

// NotEq creates a new inequality condition for the specified field and value.
//
// Parameters:
//   - field: The database column or field name
//   - value: The value to compare against the field
//
// Example:
//
//	// Creates: WHERE status != 'active'
//	b.Where(builder.NotEq("status", "active"))
func NotEq(field string, value interface{}) *Condition {
	return newCondition(false, field, "!=", []interface{}{value})
}

// Gt creates a new greater than condition for the specified field and value.
//
// Parameters:
//   - field: The database column or field name
//   - value: The value to compare against the field
//
// Example:
//
//	// Creates: WHERE age > 18
//	b.Where(builder.Gt("age", 18))
func Gt(field string, value interface{}) *Condition {
	return newCondition(false, field, ">", []interface{}{value})
}

// Gte creates a new greater than or equal to condition for the specified field and value.
//
// Parameters:
//   - field: The database column or field name
//   - value: The value to compare against the field
//
// Example:
//
//	// Creates: WHERE age >= 18
//	b.Where(builder.Gte("age", 18))
func Gte(field string, value interface{}) *Condition {
	return newCondition(false, field, ">=", []interface{}{value})
}

// Lt creates a new less than condition for the specified field and value.
//
// Parameters:
//   - field: The database column or field name
//   - value: The value to compare against the field
//
// Example:
//
//	// Creates: WHERE age < 18
//	b.Where(builder.Lt("age", 18))
func Lt(field string, value interface{}) *Condition {
	return newCondition(false, field, "<", []interface{}{value})
}

// Lte creates a new less than or equal to condition for the specified field and value.
//
// Parameters:
//   - field: The database column or field name
//   - value: The value to compare against the field
//
// Example:
//
//	// Creates: WHERE age <= 18
//	b.Where(builder.Lte("age", 18))
func Lte(field string, value interface{}) *Condition {
	return newCondition(false, field, "<=", []interface{}{value})
}

// Like creates a new LIKE condition for the specified field and value.
//
// Parameters:
//   - field: The database column or field name
//   - value: The value to compare against the field
//
// Example:
//
//	// Creates: WHERE name LIKE '%John%'
//	b.Where(builder.Like("name", "%John%"))
func Like(field string, value interface{}) *Condition {
	return newCondition(false, field, "LIKE", []interface{}{value})
}

// NotLike creates a new NOT LIKE condition for the specified field and value.
//
// Parameters:
//   - field: The database column or field name
//   - value: The value to compare against the field
//
// Example:
//
//	// Creates: WHERE name NOT LIKE '%John%'
//	b.Where(builder.NotLike("name", "%John%"))
func NotLike(field string, value interface{}) *Condition {
	return newCondition(false, field, "NOT LIKE", []interface{}{value})
}

// In creates a new IN condition for the specified field.
// The values parameter contains the list of values to match against.
//
// Parameters:
//   - field: The database column or field name
//   - values: One or more values to include in the IN clause
//
// Example:
//
//	// Creates: WHERE status IN ('pending', 'processing', 'completed')
//	b.Where(builder.In("status", "pending", "processing", "completed"))
func In(field string, values ...interface{}) *Condition {
	return newCondition(false, field, "IN", values)
}

// NotIn creates a new NOT IN condition for the specified field.
// The values parameter contains the list of values to exclude.
//
// Parameters:
//   - field: The database column or field name
//   - values: One or more values to exclude in the NOT IN clause
//
// Example:
//
//	// Creates: WHERE status NOT IN ('deleted', 'banned')
//	b.Where(builder.NotIn("status", "deleted", "banned"))
func NotIn(field string, values ...interface{}) *Condition {
	return newCondition(false, field, "NOT IN", values)
}

// NewConditionGroup creates a group of conditions that can be used together.
// It accepts multiple conditions and returns them as a slice.
//
// Warning: Do not mix OrderBy and Where conditions in the same group as they
// serve different purposes and may lead to unexpected behavior.
func NewConditionGroup(conds ...*Condition) (cg []*Condition) {
	if len(conds) > 0 {
		cg = append(cg, conds...)
	}
	return
}

// // NewConditionAnd ...
// func NewConditionAnd(field string, op string, values []interface{}) Condition {
// 	return newCondition(true, field, op, values)
// }

// // NewConditionOr ...
// func NewConditionOr(field string, op string, values []interface{}) Condition {
// 	return newCondition(false, field, op, values)
// }

// newOrderCondition creates a new Condition for ORDER BY clauses.
// It is an internal helper function used by Asc and Desc functions.
func newOrderCondition(field string, asc bool) *Condition {
	return &Condition{
		Field: field,
		Asc:   asc,
	}
}

// OrderBy creates a new slice of ordering conditions.
// It accepts multiple order conditions and combines them into a single ordering clause.
func OrderBy(conds ...*Condition) (by []*Condition) {
	if len(conds) > 0 {
		by = append(by, conds...)
	}
	return
}

// Desc creates a new descending ORDER BY condition for the specified field.
// It is used to sort results in descending order.
func Desc(field string) *Condition {
	return newOrderCondition(field, false)
}

// Asc creates a new ascending ORDER BY condition for the specified field.
// It is used to sort results in ascending order.
func Asc(field string) *Condition {
	return newOrderCondition(field, true)
}
