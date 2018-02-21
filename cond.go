package builder

// type Fields []string

// Values ...
type Values []interface{}

// Condition for Where & Order Cluse
type Condition struct {
	Field    string        // field name
	Asc      bool          // true for `ASC`, false for `DESC`
	AndOr    bool          // true for `AND`, false for `OR`
	Operator string        // where operator
	Values   []interface{} // query args
}

// newCondition : get new Condition
func newCondition(andOr bool, field string, op string, values []interface{}) Condition {
	return Condition{
		AndOr:    andOr,
		Field:    field,
		Operator: op,
		Values:   values,
	}
}

// NewConditionAnd ...
func NewConditionAnd(field string, op string, values []interface{}) Condition {
	return newCondition(true, field, op, values)
}

// NewConditionOr ...
func NewConditionOr(field string, op string, values []interface{}) Condition {
	return newCondition(false, field, op, values)
}

// newOrderCondition : get new Order Condition
func newOrderCondition(field string, asc bool) Condition {
	return Condition{
		Field: field,
		Asc:   asc,
	}
}

// OrderDesc ...
func OrderDesc(field string) Condition {
	return newOrderCondition(field, false)
}

// OrderAsc ...
func OrderAsc(field string) Condition {
	return newOrderCondition(field, true)
}
