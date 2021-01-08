package builder

// SQLType ...
type SQLType int

const (
	// RawSQL ...
	RawSQL SQLType = iota + 1
	// SelectSQL ...
	SelectSQL
	// InsertSQL ...
	InsertSQL
	// ReplaceSQL ...
	ReplaceSQL
	// InsertOrUpdateSQL ...
	InsertOrUpdateSQL
	// UpdateSQL ...
	UpdateSQL
	// DeleteSQL ...
	DeleteSQL
)

// func (f *FieldValue) String() string {
// 	return f.Name
// }

// operator map: number of operator values. 3 means multi-values.
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
