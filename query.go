package builder

import (
	"fmt"
	"strings"
)

// Query ...
type Query struct {
	Query string
	Args  []interface{}
}

// NewQuery ...
func NewQuery(q string, args ...interface{}) *Query {
	return &Query{
		Query: q,
		Args:  args,
	}
}

func (q *Query) String() (str string) {
	str = strings.Replace(q.Query, " = ?", " = '%v'", -1)
	str = strings.Replace(str, " != ?", " != '%v'", -1)
	str = strings.Replace(str, " <> ?", " <> '%v'", -1)
	str = strings.Replace(str, " > ?", " > '%v'", -1)
	str = strings.Replace(str, " >= ?", " >= '%v'", -1)
	str = strings.Replace(str, " < ?", " < '%v'", -1)
	str = strings.Replace(str, " <= ?", " <= '%v'", -1)
	str = strings.Replace(str, " (?, ", " ('%v', ", -1)
	str = strings.Replace(str, " (?)", " ('%v')", -1)
	// str = strings.Replace(str, ", ?)", ", '%v')", -1)
	str = strings.Replace(str, ", ?", ", '%v'", -1)
	// str = strings.Replace(str, ", ? ", ", '%v' ", -1)
	str = fmt.Sprintf(str, q.Args...)
	return
}
