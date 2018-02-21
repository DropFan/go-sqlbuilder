package builder

import (
	"reflect"
	"testing"
)

func TestNewQuery(t *testing.T) {
	qCase1 := &Query{
		Query: "test",
		Args:  nil,
	}
	type args struct {
		q    string
		args []interface{}
	}
	tests := []struct {
		name string
		args args
		want *Query
	}{
		// TODO: Add test cases.
		{
			name: "case1",
			args: args{
				q:    "test",
				args: nil,
			},
			want: qCase1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQuery(tt.args.q, tt.args.args...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
