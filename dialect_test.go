package builder

import (
	"testing"
)

func TestMysqlDialector_Escape(t *testing.T) {
	type fields struct {
		Dialector Dialector
	}
	type args struct {
		s []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{name: "single_field", args: args{[]string{"user"}}, want: "`user`"},
		{name: "multiple_fields", args: args{[]string{"user", "age", "sex"}}, want: "`user`, `age`, `sex`"},
		{name: "empty_field", args: args{[]string{""}}, want: "``"},
		{name: "special_chars", args: args{[]string{"user.name", "table-1", "column_2"}}, want: "`user.name`, `table-1`, `column_2`"},
		{name: "with_spaces", args: args{[]string{"first name", "last name"}}, want: "`first name`, `last name`"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MysqlDialector{
				Dialector: tt.fields.Dialector,
			}
			if got := m.Escape(tt.args.s...); got != tt.want {
				t.Errorf("MysqlDialector.Escape() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresqlDialector_Escape(t *testing.T) {
	type fields struct {
		Dialector Dialector
	}
	type args struct {
		s []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{name: "single_field", args: args{[]string{"user"}}, want: `"user"`},
		{name: "multiple_fields", args: args{[]string{"user", "age", "sex"}}, want: `"user", "age", "sex"`},
		{name: "empty_field", args: args{[]string{""}}, want: `""`},
		{name: "special_chars", args: args{[]string{"user.name", "table-1", "column_2"}}, want: `"user.name", "table-1", "column_2"`},
		{name: "with_spaces", args: args{[]string{"first name", "last name"}}, want: `"first name", "last name"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PostgresqlDialector{
				Dialector: tt.fields.Dialector,
			}
			if got := p.Escape(tt.args.s...); got != tt.want {
				t.Errorf("PostgresqlDialector.Escape() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMysqlDialector_Placeholder(t *testing.T) {
	type fields struct {
		Dialector Dialector
	}
	type args struct {
		index int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{name: "positive_index", fields: fields{}, args: args{1}, want: "?"},
		{name: "zero_index", fields: fields{}, args: args{0}, want: "?"},
		{name: "large_index", fields: fields{}, args: args{999}, want: "?"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MysqlDialector{
				Dialector: tt.fields.Dialector,
			}
			if got := m.Placeholder(tt.args.index); got != tt.want {
				t.Errorf("MysqlDialector.Placeholder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresqlDialector_Placeholder(t *testing.T) {
	type fields struct {
		Dialector Dialector
	}
	type args struct {
		index int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{name: "positive_index", fields: fields{}, args: args{1}, want: "$1"},
		{name: "zero_index", fields: fields{}, args: args{0}, want: "$0"},
		{name: "large_index", fields: fields{}, args: args{999}, want: "$999"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PostgresqlDialector{
				Dialector: tt.fields.Dialector,
			}
			if got := p.Placeholder(tt.args.index); got != tt.want {
				t.Errorf("PostgresqlDialector.Placeholder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLiteDialector_Escape(t *testing.T) {
	type fields struct {
		Dialector Dialector
	}
	type args struct {
		s []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{name: "single_field", args: args{[]string{"user"}}, want: `"user"`},
		{name: "multiple_fields", args: args{[]string{"user", "age", "sex"}}, want: `"user", "age", "sex"`},
		{name: "empty_field", args: args{[]string{""}}, want: `""`},
		{name: "special_chars", args: args{[]string{"user.name", "table-1", "column_2"}}, want: `"user.name", "table-1", "column_2"`},
		{name: "with_spaces", args: args{[]string{"first name", "last name"}}, want: `"first name", "last name"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SQLiteDialector{
				Dialector: tt.fields.Dialector,
			}
			if got := s.Escape(tt.args.s...); got != tt.want {
				t.Errorf("SQLiteDialector.Escape() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLiteDialector_Placeholder(t *testing.T) {
	type fields struct {
		Dialector Dialector
	}
	type args struct {
		index int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{name: "positive_index", fields: fields{}, args: args{1}, want: "?"},
		{name: "zero_index", fields: fields{}, args: args{0}, want: "?"},
		{name: "large_index", fields: fields{}, args: args{999}, want: "?"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SQLiteDialector{
				Dialector: tt.fields.Dialector,
			}
			if got := s.Placeholder(tt.args.index); got != tt.want {
				t.Errorf("SQLiteDialector.Placeholder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLiteDialector_GetEscapeChar(t *testing.T) {
	type fields struct {
		Dialector Dialector
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "default", fields: fields{}, want: `"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := SQLiteDialector{
				Dialector: tt.fields.Dialector,
			}
			if got := s.GetEscapeChar(); got != tt.want {
				t.Errorf("SQLiteDialector.GetEscapeChar() = %v, want %v", got, tt.want)
			}
		})
	}
}
