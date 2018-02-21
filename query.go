package builder

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
