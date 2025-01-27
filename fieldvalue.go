// Package builder provides a fluent SQL query builder with support for multiple SQL dialects.
package builder

// // Field ...
// type Field struct {
// 	Name string
// }

// func (f *Field) String() string {
// 	return f.Name
// }

// func (f *Field) methodName()  {

// }

// // NewField return an new Field
// func NewField(name string) Field {
// 	return Field{
// 		Name: name,
// 	}
// }

// FieldValue represents a field-value pair used in SQL queries for setting values
// in INSERT, UPDATE, and other operations. It provides a convenient way to associate
// a field name with its corresponding value while maintaining type safety.
type FieldValue struct {
	// Name is the name of the database field or column
	Name string

	// Value is the value to be assigned to the field, can be of any type
	Value interface{}
}

// NewFieldValue creates a new FieldValue instance with the specified field name and value.
// It is the primary constructor for creating field-value pairs.
//
// Parameters:
//   - name: The name of the database field or column
//   - value: The value to be assigned to the field
// 
// Returns:
//   - *FieldValue: A pointer to the newly created FieldValue instance
//
// Example:
//
//	fv := NewFieldValue("age", 25)
//	b.Update("users").Set(fv)
func NewFieldValue(name string, value interface{}) *FieldValue {
	return &FieldValue{
		Name:  name,
		Value: value,
	}
}

// NewFV is a shorthand alias for NewFieldValue. It creates a new FieldValue instance
// with the specified field name and value.
//
// Parameters:
//   - name: The name of the database field or column
//   - value: The value to be assigned to the field
//
// Returns:
//   - *FieldValue: A pointer to the newly created FieldValue instance
//
// Example:
//
//	fv := NewFV("status", "active")
func NewFV(name string, value interface{}) *FieldValue {
	return NewFieldValue(name, value)
}

// NewKV is a shorthand alias for NewFieldValue. It creates a new FieldValue instance
// with the specified field name and value. The name 'KV' stands for 'Key-Value'.
//
// Parameters:
//   - name: The name of the database field or column
//   - value: The value to be assigned to the field
//
// Returns:
//   - *FieldValue: A pointer to the newly created FieldValue instance
//
// Example:
//
//	fv := NewKV("email", "user@example.com")
func NewKV(name string, value interface{}) *FieldValue {
	return NewFieldValue(name, value)
}
