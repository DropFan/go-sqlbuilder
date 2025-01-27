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

func TestIn(t *testing.T) {
	type args struct {
		field  string
		values []interface{}
	}

	tests := []struct {
		name string
		args args
		want *Condition
	}{
		// TODO: Add test cases.
		{
			name: "in",
			args: args{
				field:  "name",
				values: []interface{}{"coder", "hacker"},
			},
			want: &Condition{
				Field:    "name",
				Operator: "IN",
				Values:   []interface{}{"coder", "hacker"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCond := In(tt.args.field, tt.args.values...); !reflect.DeepEqual(gotCond, tt.want) {
				t.Errorf("In() = \n%#v\n, want\n%#v", gotCond, tt.want)
			}
		})
	}
}

func TestNotIn(t *testing.T) {
	type args struct {
		field  string
		values []interface{}
	}

	tests := []struct {
		name string
		args args
		want *Condition
	}{
		// TODO: Add test cases.
		{
			name: "not in",
			args: args{
				field:  "name",
				values: []interface{}{"coder", "hacker"},
			},
			want: &Condition{
				Field:    "name",
				Operator: "NOT IN",
				Values:   []interface{}{"coder", "hacker"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCond := NotIn(tt.args.field, tt.args.values...); !reflect.DeepEqual(gotCond, tt.want) {
				t.Errorf("NotIn() = \n%#v\n, want\n%#v", gotCond, tt.want)
			}
		})
	}
}

func TestBetween(t *testing.T) {
	type args struct {
		field  string
		values []interface{}
	}

	tests := []struct {
		name string
		args args
		want *Condition
	}{
		// TODO: Add test cases.
		{
			name: "between",
			args: args{
				field:  "age",
				values: []interface{}{12, 36},
			},
			want: &Condition{
				Field:    "age",
				Operator: "BETWEEN",
				Values:   []interface{}{12, 36},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCond := Between(tt.args.field, tt.args.values...); !reflect.DeepEqual(gotCond, tt.want) {
				t.Errorf("Between() = \n%#v\n, want\n%#v", gotCond, tt.want)
			}
		})
	}
}

func TestNotBetween(t *testing.T) {
	type args struct {
		field  string
		values []interface{}
	}

	tests := []struct {
		name string
		args args
		want *Condition
	}{
		// TODO: Add test cases.
		{
			name: "not between",
			args: args{
				field:  "age",
				values: []interface{}{12, 36},
			},
			want: &Condition{
				Field:    "age",
				Operator: "NOT BETWEEN",
				Values:   []interface{}{12, 36},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCond := NotBetween(tt.args.field, tt.args.values...); !reflect.DeepEqual(gotCond, tt.want) {
				t.Errorf("NotBetween() = \n%#v\n, want\n%#v", gotCond, tt.want)
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

func TestEq(t *testing.T) {
	type args struct {
		field string
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want *Condition
	}{
		{
			name: "status=active",
			args: args{
				field: "status",
				value: "active",
			},
			want: &Condition{
				Field:    "status",
				Operator: "=",
				Values:   []interface{}{"active"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCond := Eq(tt.args.field, tt.args.value); !reflect.DeepEqual(gotCond, tt.want) {
				t.Errorf("Eq() = \n%#v\n, want\n%#v", gotCond, tt.want)
			}
		})
	}
}

func TestNotEq(t *testing.T) {
	type args struct {
		field string
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want *Condition
	}{
		{
			name: "status!=inactive",
			args: args{
				field: "status",
				value: "inactive",
			},
			want: &Condition{
				Field:    "status",
				Operator: "!=",
				Values:   []interface{}{"inactive"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCond := NotEq(tt.args.field, tt.args.value); !reflect.DeepEqual(gotCond, tt.want) {
				t.Errorf("NotEq() = \n%#v\n, want\n%#v", gotCond, tt.want)
			}
		})
	}
}

func TestGt(t *testing.T) {
	type args struct {
		field string
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want *Condition
	}{
		{
			name: "age>18",
			args: args{
				field: "age",
				value: 18,
			},
			want: &Condition{
				Field:    "age",
				Operator: ">",
				Values:   []interface{}{18},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCond := Gt(tt.args.field, tt.args.value); !reflect.DeepEqual(gotCond, tt.want) {
				t.Errorf("Gt() = \n%#v\n, want\n%#v", gotCond, tt.want)
			}
		})
	}
}

func TestGte(t *testing.T) {
	type args struct {
		field string
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want *Condition
	}{
		{
			name: "age>=18",
			args: args{
				field: "age",
				value: 18,
			},
			want: &Condition{
				Field:    "age",
				Operator: ">=",
				Values:   []interface{}{18},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCond := Gte(tt.args.field, tt.args.value); !reflect.DeepEqual(gotCond, tt.want) {
				t.Errorf("Gte() = \n%#v\n, want\n%#v", gotCond, tt.want)
			}
		})
	}
}

func TestLt(t *testing.T) {
	type args struct {
		field string
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want *Condition
	}{
		{
			name: "age<18",
			args: args{
				field: "age",
				value: 18,
			},
			want: &Condition{
				Field:    "age",
				Operator: "<",
				Values:   []interface{}{18},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCond := Lt(tt.args.field, tt.args.value); !reflect.DeepEqual(gotCond, tt.want) {
				t.Errorf("Lt() = \n%#v\n, want\n%#v", gotCond, tt.want)
			}
		})
	}
}

func TestLte(t *testing.T) {
	type args struct {
		field string
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want *Condition
	}{
		{
			name: "age<=18",
			args: args{
				field: "age",
				value: 18,
			},
			want: &Condition{
				Field:    "age",
				Operator: "<=",
				Values:   []interface{}{18},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCond := Lte(tt.args.field, tt.args.value); !reflect.DeepEqual(gotCond, tt.want) {
				t.Errorf("Lte() = \n%#v\n, want\n%#v", gotCond, tt.want)
			}
		})
	}
}

func TestLike(t *testing.T) {
	type args struct {
		field string
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want *Condition
	}{
		{
			name: "name like John",
			args: args{
				field: "name",
				value: "%John%",
			},
			want: &Condition{
				Field:    "name",
				Operator: "LIKE",
				Values:   []interface{}{"%John%"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCond := Like(tt.args.field, tt.args.value); !reflect.DeepEqual(gotCond, tt.want) {
				t.Errorf("Like() = \n%#v\n, want\n%#v", gotCond, tt.want)
			}
		})
	}
}

func TestNotLike(t *testing.T) {
	type args struct {
		field string
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want *Condition
	}{
		{
			name: "name not like John",
			args: args{
				field: "name",
				value: "%John%",
			},
			want: &Condition{
				Field:    "name",
				Operator: "NOT LIKE",
				Values:   []interface{}{"%John%"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCond := NotLike(tt.args.field, tt.args.value); !reflect.DeepEqual(gotCond, tt.want) {
				t.Errorf("NotLike() = \n%#v\n, want\n%#v", gotCond, tt.want)
			}
		})
	}
}
