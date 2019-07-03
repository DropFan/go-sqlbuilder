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
		// TODO: Add test cases.
		{name: "test_1", args: args{[]string{"user"}}, want: "`user`"},
		{name: "test_2", args: args{[]string{"user", "age", "sex"}}, want: "`user`, `age`, `sex`"},
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
		// TODO: Add test cases.
		{name: "test_1", args: args{[]string{"user"}}, want: `"user"`},
		{name: "test_2", args: args{[]string{"user", "age", "sex"}}, want: `"user", "age", "sex"`},
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
		{"test_1", fields{}, args{11}, "?"},
		{"test_2", fields{}, args{1}, "?"},
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
		{"test_1", fields{}, args{1}, "$1"},
		{"test_2", fields{}, args{2}, "$2"},
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
