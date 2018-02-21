package builder

// SQLType ...
type SQLType int

const (
	// RawSQL ...
	RawSQL = iota + 1
	// SelectSQL ...
	SelectSQL
	// InsertSQL ...
	InsertSQL
	// UpdateSQL ...
	UpdateSQL
	// DeleteSQL ...
	DeleteSQL
)

// FieldValue ...
type FieldValue struct {
	Name  string
	Value interface{}
}

// func (f *FieldValue) String() string {
// 	return f.Name
// }

// operator map: number of operator values
var operMap = map[string]int{
	"=":           1,
	"!=":          1,
	"<>":          1,
	">":           1,
	">=":          1,
	"<":           1,
	"<=":          1,
	"IN":          3,
	"NOT IN":      3,
	"LIKE":        1,
	"NOT LIKE":    1,
	"BETWEEN":     2,
	"NOT BETWEEN": 2,
}
