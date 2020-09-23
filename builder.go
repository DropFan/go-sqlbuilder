package builder

import (
	"fmt"
	"strconv"
	"strings"
)

// Builder :You'd better create a new builder with `New()` function
//
// You can found some examples in `builder_test.go`
type Builder struct {
	sqlType     SQLType
	dialector   Dialector
	queryTables string // abandoned
	queryArgs   []interface{}
	query       string
	from        string // abandoned
	setValues   []string
	where       []string // abandoned
	orderBy     string   // abandoned
	limit       string   // abandoned
	ErrList     []error
	lastQueries []*Query
	// placeholder
}

// SetDialector set dialector for params binding placeholder and escape charcter
func (b *Builder) SetDialector(d Dialector) *Builder {
	b.dialector = d
	return b
}

// Escape escape fields
func (b *Builder) Escape(s ...string) string {
	return b.dialector.Escape(s...)
}

// EscapeChar return escape char
func (b *Builder) EscapeChar() string {
	return b.dialector.GetEscapeChar()
}

// Placeholder ...
// func (b *Builder) Placeholder(index int) string {
// 	return b.dialector.Placeholder(index)
// }

// New builder
// You can found some examples in `builder_test.go`
func New() *Builder {
	return &Builder{
		sqlType:     0,
		dialector:   mysqlDialector,
		queryTables: "",
		queryArgs:   []interface{}{},
		query:       "",
		from:        "",
		setValues:   []string{},
		where:       []string{},
		orderBy:     "",
		limit:       "",
		ErrList:     []error{},
		lastQueries: []*Query{},
	}
}

// LastQueries return last queries
func (b *Builder) LastQueries() []*Query {
	return b.lastQueries
}

// LastQuery return last query
func (b *Builder) LastQuery() *Query {
	if len(b.lastQueries) == 0 {
		return nil
	}

	return b.lastQueries[len(b.lastQueries)-1]
}

// renew reset some data after `Build()` was called
func (b *Builder) renew(st SQLType) {
	b.ErrList = []error{}
	b.sqlType = st
	b.queryArgs = []interface{}{}
	b.query = ""
	b.setValues = []string{}
	b.where = []string{}
	b.orderBy = ""
	b.limit = ""
}

// Clear clear current query & query args
func (b *Builder) Clear() *Builder {
	b.renew(0)
	return b
}

// QueryArgs return current query args
func (b *Builder) QueryArgs() []interface{} {
	return b.queryArgs
}

// Query return current sql
func (b *Builder) Query() string {
	return b.query
}

// Append append query & query args to origin query
func (b *Builder) Append(s string, args ...interface{}) *Builder {
	b.query += s
	b.queryArgs = append(b.queryArgs, args...)
	return b
}

// AppendPre append query & query args to ahead of origin query
func (b *Builder) AppendPre(s string, args ...interface{}) *Builder {
	b.query = s + b.query
	b.queryArgs = append(args, b.queryArgs...)
	return b
}

// Raw ...
func (b *Builder) Raw(s string, args ...interface{}) *Builder {
	b.renew(RawSQL)
	b.query = s
	b.queryArgs = append(args, b.queryArgs...)
	return b
}

// Select ...
func (b *Builder) Select(fields ...string) *Builder {
	b.renew(SelectSQL)
	b.query = "SELECT"

	if len(fields) <= 0 {
		b.query += ""
	} else if fields[0] == "*" {
		b.query += " *"
	} else {
		b.query += " " + b.Escape(fields...)
		// b.query += " `" + strings.Join(fields, "`, `") + "`"
	}

	return b
}

// Insert ...
func (b *Builder) Insert(tableName string, fields ...string) *Builder {
	b.renew(InsertSQL)
	b.query = "INSERT INTO " + b.Escape(tableName)

	if len(fields) > 0 {
		b.Into(fields...)
	}
	return b
}

// Replace ...
func (b *Builder) Replace(tableName string, fields ...string) *Builder {
	b.renew(InsertSQL)
	b.query = "REPLACE INTO " + b.Escape(tableName)

	if len(fields) > 0 {
		b.Into(fields...)
	}
	return b
}

// Into ... for insert or replace
func (b *Builder) Into(fields ...string) *Builder {
	b.query += " (" + b.Escape(fields...) + ")"
	// b.query += " (`" + strings.Join(fields, "`, `") + "`)"
	return b
}

// Values ...
func (b *Builder) Values(valsGroup ...[]interface{}) *Builder {
	b.query += " VALUES "
	// index := 0
	for i, vals := range valsGroup {
		if i > 0 {
			b.query += ", "
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

		b.query += "(?" + strings.Repeat(", ?", len(vals)-1) + ")"
		b.queryArgs = append(b.queryArgs, vals...)
	}
	return b
}

// Update ...
func (b *Builder) Update(tableName string, fvals ...*FieldValue) *Builder {
	b.renew(UpdateSQL)
	b.query = "UPDATE " + b.Escape(tableName) + " SET "

	if len(fvals) > 0 {
		b.Set(fvals...)
	}

	return b
}

// Set ...
func (b *Builder) Set(fvals ...*FieldValue) *Builder {
	// b.setValue = ""

	for i, fval := range fvals {
		if fval == nil {
			continue
		}
		if i > 0 || len(b.setValues) > 0 {
			b.query += ", "
		}
		b.setValues = append(b.setValues, fval.Name)
		b.query += b.Escape(fval.Name) + " = ?"
		b.queryArgs = append(b.queryArgs, fval.Value)
	}

	return b
}

// Delete ...
func (b *Builder) Delete(tableName string) *Builder {
	b.renew(DeleteSQL)
	b.query = "DELETE FROM " + b.Escape(tableName)

	return b
}

// Build ...
func (b *Builder) Build(queries ...interface{}) (q *Query, err error) {

	switch b.sqlType {
	case SelectSQL:
	case InsertSQL:
	case UpdateSQL:
	case DeleteSQL:
	case RawSQL:
	default:
		return nil, ErrEmptySQLType
	}
	if len(b.ErrList) > 0 {
		err = ErrListIsNotEmpty
	}
	q = NewQuery(b.query, b.queryArgs...)
	b.lastQueries = append(b.lastQueries, q)
	b.renew(RawSQL)
	return q, err
}

// From ...
func (b *Builder) From(tables ...string) *Builder {
	if len(tables) <= 0 {
		return b
	}
	// b.Tables = tables
	// b.QueryTables = "`" + strings.Join(tables, "`, `") + "`"
	b.query += " FROM " + b.Escape(tables...)
	// b.query += " FROM `" + strings.Join(tables, "`, `") + "`"
	return b
}

// FromRaw ...
func (b *Builder) FromRaw(from string) *Builder {
	b.query += " FROM " + from
	return b
}

// Where ...
func (b *Builder) addConditions(conditions ...*Condition) *Builder {
	condSlice := []string{}
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
	b.query += strings.Join(condSlice, " ")
	return b

}

// And ...
func (b *Builder) And(conditions ...*Condition) *Builder {
	switch len(conditions) {
	case 0:
		return b
	case 1:
		b.query += " AND "
		b.addConditions(conditions...)
	default:
		b.query += " AND ("
		b.addConditions(conditions...)
		b.query += ")"
	}
	// if len(b.where) > 0 {
	// 	b.where[0] = "(" + b.where[0]
	// 	b.where[len(b.where)-1] = b.where[len(b.where)-1] + ") AND"
	// }
	// b.where = append(b.where, ") AND")
	// b.where = append(b.where, ")")
	return b
}

// Or ...
func (b *Builder) Or(conditions ...*Condition) *Builder {
	switch len(conditions) {
	case 0:
		return b
	case 1:
		b.query += " OR "
		b.addConditions(conditions...)
	default:
		b.query += " OR ("
		b.addConditions(conditions...)
		b.query += ")"
	}

	// if len(b.where) > 0 {
	// 	b.where[0] = "(" + b.where[0]
	// 	b.where[len(b.where)-1] = b.where[len(b.where)-1] + ") OR"
	// }
	// b.addConditions(conditions...)
	// b.where = append(b.where, "")
	return b
}

// In ...
func (b *Builder) In(field string, values ...interface{}) *Builder {
	b.addConditions(In(field, values...))
	return b
}

// NotIn ...
func (b *Builder) NotIn(field string, values ...interface{}) *Builder {
	b.addConditions(NotIn(field, values...))
	return b
}

// Between ...
func (b *Builder) Between(field string, values ...interface{}) *Builder {
	b.addConditions(Between(field, values...))
	return b
}

// NotBetween ...
func (b *Builder) NotBetween(field string, values ...interface{}) *Builder {
	b.addConditions(NotBetween(field, values...))
	return b
}

// Where ...
func (b *Builder) Where(conditions ...*Condition) *Builder {
	b.query += " WHERE "
	if len(conditions) == 0 {
		b.query += "1"
		return b
	}

	b.addConditions(conditions...)
	return b
}

// WhereRaw ...
func (b *Builder) WhereRaw(str string, args ...interface{}) *Builder {
	b.query += " WHERE "
	b.query += str
	b.queryArgs = append(b.queryArgs, args...)

	return b
}

// BuildWhere ...
func (b *Builder) buildCondition(cond *Condition) (str string, queryArgs []interface{}, err error) {
	// if cond == nil {
	// 	return "", nil, ErrEmptyCondition
	// }

	str = ""
	queryArgs = []interface{}{}

	if opValue, ok := operMap[cond.Operator]; !ok {
		// return "", queryArgs,
		err = fmt.Errorf("Invalid operator:(operator:%s[field:%s])", cond.Operator, cond.Field)
		return
	} else if len(cond.Values) != opValue {
		switch opValue {
		case 3:
			if len(cond.Values) != 0 {
				break
			}
			fallthrough
		case 1, 2:
			err = fmt.Errorf("Invalid number of values with operator:(%s[field:%s])", cond.Operator, cond.Field)
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

// OrderBy ...
func (b *Builder) OrderBy(conditions ...*Condition) *Builder {
	// order, err := buildOrderBy(conditions)
	b.query += " ORDER BY "

	condStrSlice := []string{}

	for _, cond := range conditions {
		if cond == nil {
			continue
		}
		condStr := b.Escape(cond.Field)
		if cond.Asc {
			condStr += " ASC"
		} else {
			condStr += " DESC"
		}
		condStrSlice = append(condStrSlice, condStr)
	}
	b.query += strings.Join(condStrSlice, ", ")

	// b.orderBy = b.buildOrderBy(conditions...)
	return b
}

// Limit ...
func (b *Builder) Limit(limitOffset ...int) *Builder {
	// var limit, offset string
	if len(limitOffset) == 1 {
		b.query += " LIMIT " + strconv.Itoa(limitOffset[0])
	} else {
		b.query += " LIMIT " +
			strconv.Itoa(limitOffset[0]) +
			", " +
			strconv.Itoa(limitOffset[1])
	}
	// b.queryArgs = append(b.queryArgs, offset, limit)
	return b
}

// Count ...
func (b *Builder) Count(query ...string) *Builder {
	b.renew(SelectSQL)
	b.query = "SELECT COUNT("
	if len(query) <= 0 {
		b.query += "1"
	} else {
		b.query += strings.TrimSpace(strings.Join(query, " "))
	}
	b.query += ")"

	return b
}
