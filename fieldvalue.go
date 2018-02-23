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

// FieldValue ...
type FieldValue struct {
	Name  string
	Value interface{}
}

// NewFieldValue ...
func NewFieldValue(name string, value interface{}) *FieldValue {
	return &FieldValue{
		Name:  name,
		Value: value,
	}
}

// NewFV ...
func NewFV(name string, value interface{}) *FieldValue {
	return NewFieldValue(name, value)
}

// NewKV ...
func NewKV(name string, value interface{}) *FieldValue {
	return NewFieldValue(name, value)
}
