package builder

import (
	"reflect"
	"testing"
)

var (
	ageGT1 = &Condition{
		Field:    "age",
		Operator: ">=",
		Values:   []interface{}{1},
	}
	nameEqCoder = &Condition{
		// AndOr:    true,
		Field:    "name",
		Operator: "=",
		Values:   []interface{}{"coder"},
	}
	nameInNames = &Condition{
		AndOr:    true,
		Field:    "name",
		Operator: "IN",
		Values:   []interface{}{"coder", "hacker"},
	}
	sexEqFemale = &Condition{
		Field:    "sex",
		Operator: "=",
		Values:   []interface{}{"female"},
	}
	AndSexEqFemale = &Condition{
		AndOr:    true,
		Field:    "sex",
		Values:   []interface{}{"female"},
		Operator: "=",
	}
	ageBetweenCond = &Condition{
		Field:    "age",
		Operator: "BETWEEN",
		Values:   []interface{}{12, 36},
	}

	whereConds = []*Condition{}
	// orConds := []Condition{sexEqFemale, ageGT1}

	errOpCond = &Condition{
		Field:    "test",
		Operator: "!",
	}
	errValNumCond = &Condition{
		Field:    "test_field",
		Operator: "=",
		Values:   []interface{}{1, 2},
	}
	ageDesc = &Condition{
		Field: "age",
		Asc:   false,
	}
	nameAsc = &Condition{
		Field: "name",
		Asc:   true,
	}
)

func Test_newCondition(t *testing.T) {
	type args struct {
		andOr  bool
		field  string
		op     string
		values []interface{}
	}
	tests := []struct {
		name string
		args args
		want *Condition
	}{
		// TODO: Add test cases.
		{
			name: "or name=hacker",
			args: args{
				field:  "name",
				op:     "=",
				values: []interface{}{"hacker"},
			},
			want: &Condition{
				AndOr:    false,
				Field:    "name",
				Operator: "=",
				Values:   []interface{}{"hacker"},
			},
		},
		{
			name: "and name=hacker",
			args: args{
				andOr:  true,
				field:  "name",
				op:     "=",
				values: []interface{}{"hacker"},
			},
			want: &Condition{
				AndOr:    true,
				Field:    "name",
				Operator: "=",
				Values:   []interface{}{"hacker"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCond := newCondition(tt.args.andOr, tt.args.field, tt.args.op, tt.args.values); !reflect.DeepEqual(gotCond, tt.want) {
				t.Errorf("newCondition() = %v, want %v", gotCond, tt.want)
			}
		})
	}
}

func TestAnd(t *testing.T) {
	type args struct {
		field  string
		op     string
		values []interface{}
	}
	tests := []struct {
		name string
		args args
		want *Condition
	}{
		// TODO: Add test cases.
		{
			name: "name=hacker",
			args: args{
				field:  "name",
				op:     "=",
				values: []interface{}{"hacker"},
			},
			want: &Condition{
				AndOr:    true,
				Field:    "name",
				Operator: "=",
				Values:   []interface{}{"hacker"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCond := And(tt.args.field, tt.args.op, tt.args.values...); !reflect.DeepEqual(gotCond, tt.want) {
				t.Errorf("And() = \n%#v\n, want\n%#v", gotCond, tt.want)
			}
		})
	}
}

func TestOr(t *testing.T) {
	type args struct {
		field  string
		op     string
		values []interface{}
	}
	tests := []struct {
		name string
		args args
		want *Condition
	}{
		// TODO: Add test cases.
		{
			name: "name=hacker",
			args: args{
				field:  "name",
				op:     "=",
				values: []interface{}{"hacker"},
			},
			want: &Condition{
				AndOr:    false,
				Field:    "name",
				Operator: "=",
				Values:   []interface{}{"hacker"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCond := Or(tt.args.field, tt.args.op, tt.args.values...); !reflect.DeepEqual(gotCond, tt.want) {
				t.Errorf("Or() = \n%#v\n, want\n%#v", gotCond, tt.want)
			}
		})
	}
}

func Test_newOrderCondition(t *testing.T) {
	type args struct {
		field string
		asc   bool
	}
	tests := []struct {
		name string
		args args
		want *Condition
	}{
		// TODO: Add test cases.
		{
			name: "by age desc",
			args: args{
				field: "age",
				asc:   false,
			},
			want: &Condition{
				Field: "age",
				Asc:   false,
			},
		},
		{
			name: "by age asc",
			args: args{
				field: "age",
				asc:   true,
			},
			want: &Condition{
				Field: "age",
				Asc:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCond := newOrderCondition(tt.args.field, tt.args.asc); !reflect.DeepEqual(gotCond, tt.want) {
				t.Errorf("newOrderCondition() = %+v, want %+v", gotCond, tt.want)
			}
		})
	}
}

func TestDesc(t *testing.T) {
	type args struct {
		field string
	}
	tests := []struct {
		name string
		args args
		want *Condition
	}{
		// TODO: Add test cases.
		{
			name: "by age desc",
			args: args{
				field: "age",
				// asc:   false,
			},
			want: &Condition{
				Field: "age",
				Asc:   false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Desc(tt.args.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderDesc() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestAsc(t *testing.T) {
	type args struct {
		field string
	}
	tests := []struct {
		name string
		args args
		want *Condition
	}{
		// TODO: Add test cases.
		{
			name: "by age asc",
			args: args{
				field: "age",
				// asc:   false,
			},
			want: &Condition{
				Field: "age",
				Asc:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Asc(tt.args.field); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderAsc() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestNewConditionGroup(t *testing.T) {
	tests := []struct {
		name  string
		conds []*Condition
		want  []*Condition
	}{
		// TODO: Add test cases.
		{
			name: "test_condition_group",
			conds: []*Condition{
				nameEqCoder,
				ageBetweenCond,
			},
			want: []*Condition{
				nameEqCoder,
				ageBetweenCond,
			},
		},
		{
			name:  "test_condition_group",
			conds: []*Condition{
				// nameEqCoder,
				// ageBetweenCond,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConditionGroup(tt.conds...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConditionGroup() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestOrderBy(t *testing.T) {
	tests := []struct {
		name  string
		conds []*Condition
		want  []*Condition
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			conds: []*Condition{
				ageDesc,
				nameAsc,
			},
			want: []*Condition{
				ageDesc,
				nameAsc,
			},
		},
		{
			name:  "test2",
			conds: []*Condition{
				// nameAsc,
				// ageDesc,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OrderBy(tt.conds...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderBy() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
