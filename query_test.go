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

func TestQuery_String(t *testing.T) {
	tests := []struct {
		name string
		q    *Query
		sql  string
		args []interface{}
		want string
	}{
		{
			name: "test 1",
			sql:  "SELECT * FROM `tablename` WHERE `k` = ?",
			args: []interface{}{12135183350725},
			want: "SELECT * FROM `tablename` WHERE `k` = '12135183350725'",
		},
		{
			name: "test 2",
			sql:  "SELECT * FROM `tablename` WHERE `k` = ? AND `k2` = ?",
			args: []interface{}{12135183350725, "testv"},
			want: "SELECT * FROM `tablename` WHERE `k` = '12135183350725' AND `k2` = 'testv'",
		},
		{
			name: "test 3",
			sql:  "UPDATE `tablename` SET `f` = ? WHERE `k` = ? AND `k2` = ?",
			args: []interface{}{"fff", 12135183350725, "testv"},
			want: "UPDATE `tablename` SET `f` = 'fff' WHERE `k` = '12135183350725' AND `k2` = 'testv'",
		},
		{
			name: "test 4",
			sql:  "UPDATE `tablename` SET `f` = ?, `g` = ? WHERE `k` = ? AND `k2` = ?",
			args: []interface{}{"fff", "ggg", 12135183350725, "testv"},
			want: "UPDATE `tablename` SET `f` = 'fff', `g` = 'ggg' WHERE `k` = '12135183350725' AND `k2` = 'testv'",
		},
		{
			name: "test 5",
			sql:  "INSERT INTO `tablename` VALUES (?, ?, ?, ?)",
			args: []interface{}{"fff", "ggg", 12135183350725, "testv"},
			want: "INSERT INTO `tablename` VALUES ('fff', 'ggg', '12135183350725', 'testv')",
		},
		// New test cases for enhanced functionality
		{
			name: "PostgreSQL placeholders",
			sql:  `SELECT * FROM "users" WHERE "name" = $1 AND "age" > $2`,
			args: []interface{}{"john", 18},
			want: `SELECT * FROM "users" WHERE "name" = 'john' AND "age" > '18'`,
		},
		{
			name: "PostgreSQL IN clause",
			sql:  `SELECT * FROM "users" WHERE "id" IN ($1, $2, $3)`,
			args: []interface{}{1, 2, 3},
			want: `SELECT * FROM "users" WHERE "id" IN ('1', '2', '3')`,
		},
		{
			name: "single quote escaping",
			sql:  "SELECT * FROM users WHERE name = ?",
			args: []interface{}{"O'Reilly"},
			want: "SELECT * FROM users WHERE name = 'O''Reilly'",
		},
		{
			name: "multiple single quotes",
			sql:  "INSERT INTO comments (text) VALUES (?)",
			args: []interface{}{"It's John's book, isn't it?"},
			want: "INSERT INTO comments (text) VALUES ('It''s John''s book, isn''t it?')",
		},
		{
			name: "NULL value",
			sql:  "UPDATE users SET deleted_at = ? WHERE id = ?",
			args: []interface{}{nil, 1},
			want: "UPDATE users SET deleted_at = NULL WHERE id = '1'",
		},
		{
			name: "boolean values",
			sql:  "UPDATE users SET active = ?, verified = ? WHERE id = ?",
			args: []interface{}{true, false, 1},
			want: "UPDATE users SET active = TRUE, verified = FALSE WHERE id = '1'",
		},
		{
			name: "no arguments",
			sql:  "SELECT * FROM users",
			args: []interface{}{},
			want: "SELECT * FROM users",
		},
		{
			name: "percent sign in value",
			sql:  "SELECT * FROM users WHERE name LIKE ?",
			args: []interface{}{"%john%"},
			want: "SELECT * FROM users WHERE name LIKE '%john%'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q = NewQuery(tt.sql, tt.args...)

			if got := tt.q.String(); got != tt.want {
				t.Errorf("Query.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
