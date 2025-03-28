package builder

import (
	"reflect"
	"testing"
)

type Dao struct {
	b         *Builder
	tableName string
	fields    []string
}

// NewDao ...
func NewDao(tablename string, fields []string) *Dao {
	return &Dao{
		b:         &Builder{},
		tableName: "",
		fields:    fields,
	}
}

// Fields return fields
func (d *Dao) Fields() []string {
	return d.fields
}

// SetFields set Fields
func (d *Dao) SetFields(fields ...string) *Dao {
	d.fields = fields
	return d
}

// TableName return table name
func (d *Dao) TableName() string {
	return d.tableName
}

// SetTableName set table name
func (d *Dao) SetTableName(table string) *Dao {
	d.tableName = table
	return d
}

var (
	b = New()
	d = &Dao{
		b:         b,
		tableName: "user",
		fields:    []string{"id", "name", "age", "sex", "birthday"},
	}
)

func TestBuilder(t *testing.T) {
	var (
		got, want      string
		args, wantArgs []interface{}
		err            error
		q              *Query
	)

	q = b.LastQuery()
	if q != nil {
		t.Errorf("unexpected error")
	}

	want = "SELECT * FROM `user` WHERE `user_id` = ?"
	wantArgs = []interface{}{1}
	b.Select("*").FromRaw("`user`").WhereRaw("`user_id` = ?", 1)
	q, err = b.Build("show tables")
	got = q.Query
	args = q.Args
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
	if !reflect.DeepEqual(wantArgs, args) {
		t.Errorf("\ngotArgs:\n%#v\nwantArgs:\n%#v\n", args, wantArgs)
	}

	// lastQueries := b.LastQueries()
	// lastQuery := lastQueries[len(lastQueries)-1]
	lastQuery := b.LastQuery()
	if lastQuery.Query != got {
		t.Errorf("\ngot:\n%s\nlast query:\n%s\n", got, lastQuery.Query)
	}

}

func TestSelect(t *testing.T) {
	var (
		got, want      string
		args, wantArgs []interface{}
		q              *Query
		err            error
	)

	want = "SELECT `id`, `name`, `age`, `sex`, `birthday` FROM `user` WHERE `age` >= ? AND `name` IN (?, ?) AND `sex` = ? OR `age` BETWEEN ? AND ? OR `name` = ? ORDER BY `age` DESC, `name` ASC LIMIT 0, 100"
	wantArgs = []interface{}{1, "coder", "hacker", "female", 12, 36, "coder"}
	b.Select(d.Fields()...).
		From("user").
		Where(ageGT1, nameInNames).
		And(sexEqFemale).
		And().
		Or(ageBetweenCond).
		Or(nameEqCoder).
		OrderBy(ageDesc, nameAsc).
		Limit(0, 100)
	q, err = b.Build()
	got = q.Query
	args = q.Args

	if err != nil {
		t.Errorf("select error:%s\n%s", err, got)
	}
	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
	if !reflect.DeepEqual(wantArgs, args) {
		t.Errorf("\ngotArgs:\n%#v\nwantArgs:\n%#v\n", args, wantArgs)
	}
	want = "SELECT `id`, `name`, `age`, `sex`, `birthday` FROM `user` WHERE {error: invalid operator:(operator:![field:test])} AND `name` IN (?, ?) OR `sex` = ? OR `name` = ? AND {error: invalid number of values with operator:(=[field:test_field])} ORDER BY `age` DESC, `name` ASC LIMIT 100 LIMIT 0, 100"
	wantArgs = []interface{}{"coder", "hacker", "female", "coder"}
	b.Select(d.Fields()...).
		From("user").
		Where([]*Condition{errOpCond, nameInNames, sexEqFemale}...).
		Or(nameEqCoder).And(errValNumCond).
		Or().
		OrderBy(ageDesc, nameAsc, nil).Limit(100).Limit(0, 100)
	q, err = b.Build()
	got = q.Query
	args = q.Args

	if err == nil {
		t.Errorf("select error:%s\n%s", err, got)
	}

	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
	if !reflect.DeepEqual(wantArgs, args) {
		t.Errorf("\ngotArgs:\n%#v\nwantArgs:\n%#v\n", args, wantArgs)
	}
	// t.Errorf("\nsb:\n")

	want = "SELECT * FROM `user` WHERE 1 AND (`name` = ? OR `sex` = ?)"
	wantArgs = []interface{}{"coder", "female"}
	b.Select("*").From("user").Where().And(nameEqCoder, sexEqFemale)
	q, err = b.Build()
	got = q.Query
	args = q.Args
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}

	if !reflect.DeepEqual(wantArgs, args) {
		t.Errorf("\ngotArgs:\n%#v\nwantArgs:\n%#v\n", args, wantArgs)
	}

	want = "SELECT * FROM `user` WHERE `age` IN (?, ?, ?) AND `city_id` IN (?, ?, ?)"
	inCityIDs := &Condition{
		Field:    "city_id",
		Operator: "IN",
		Values:   []interface{}{1, 2, 3},
	}
	wantArgs = []interface{}{1, 2, 3}
	wantArgs = append(wantArgs, []interface{}{1, 2, 3}...)
	b.Select("*").From("user").Where(In("age", []interface{}{1, 2, 3}...)).And(inCityIDs).NotIn("a", 1, 2, 3).Between("b", 1, 5).NotBetween("c", 2, 3)
	b.Select("*").From("user").Where(In("age", []interface{}{1, 2, 3}...)).Append(" AND ").In("city_id", 1, 2, 3)
	q, err = b.Build()
	got = q.Query
	args = q.Args
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}

	if !reflect.DeepEqual(wantArgs, args) {
		t.Errorf("\ngotArgs:\n%#v\nwantArgs:\n%#v\n", args, wantArgs)
	}

	want = "SELECT * FROM `user` WHERE {error: invalid number of values with operator:(IN[field:city_id])}"
	inCityIDs = &Condition{
		Field:    "city_id",
		Operator: "IN",
		// Values:   []interface{}{1, 2, 3},
	}
	wantArgs = []interface{}{}
	b.Select("*").From("user").Where(inCityIDs)
	q, err = b.Build()
	got = q.Query
	args = q.Args
	if err != nil {
		t.Logf("\nexpected error: %s\nsql:%v", err, got)
	}

	q, err = b.Select("*").From("user").Where(nil).Build()
	got = q.Query
	args = q.Args
	if err != nil {
		t.Logf("expected error: %s", err)
	}
	t.Logf("\nexpected nil error: %v\nsql:%v", err, got)
}

func TestUpdate(t *testing.T) {
	var (
		got, want      string
		args, wantArgs []interface{}
		q              *Query
		err            error
		fvals          = []*FieldValue{
			{Name: "tag", Value: "test"},
			{Name: "desc", Value: "just 4 test"},
		}
		fv = NewFV("f", "v")
		kv = NewKV("k", "v")
	)

	want = "UPDATE `user` SET `k` = ?, `f` = ?, `some_field` = ?, `tag` = ?, `desc` = ? WHERE `name` = ? AND `sex` = ? ORDER BY `age` DESC, `name` ASC LIMIT 100"
	wantArgs = []interface{}{kv.Value, fv.Value, "some_value", "test", "just 4 test", "coder", "female"}
	b.Update("user", kv).Set(fv).Append(", `some_field` = ?", "some_value").Set(fvals...)

	b.Where(nameEqCoder, AndSexEqFemale).OrderBy(ageDesc, nameAsc).Limit(100)
	q, err = b.Build()
	got = q.Query
	args = q.Args
	if err != nil {
		t.Errorf("update error:%s", err)
	}

	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
	if !reflect.DeepEqual(wantArgs, args) {
		t.Errorf("\ngotArgs:\n%#v\nwantArgs:\n%#v\n", args, wantArgs)
	}

	//
	want = "UPDATE `user` SET `k` = ?, `f` = ?, `some_field` = ?, `tag` = ?, `desc` = ? WHERE `name` = ? AND `sex` = ?"
	wantArgs = []interface{}{kv.Value, fv.Value, "some_value", "test", "just 4 test", "coder", "female"}
	b.Update("user", kv).Set(fv, nil).Append(", `some_field` = ?", "some_value").Set(fvals...)

	b.Where(nameEqCoder, AndSexEqFemale)
	q, err = b.Build()
	got = q.Query
	args = q.Args
	if err != nil {
		t.Errorf("update error:%s", err)
	}

	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
	if !reflect.DeepEqual(wantArgs, args) {
		t.Errorf("\ngotArgs:\n%#v\nwantArgs:\n%#v\n", args, wantArgs)
	}
}

func TestDelete(t *testing.T) {

	var (
		got, want      string
		args, wantArgs []interface{}
		q              *Query
		err            error
	)

	want = "DELETE FROM `user` WHERE `name` = ? AND `sex` = ? /*===*/  OR (`sex` = ? AND `name` IN (?, ?)) ORDER BY `age` DESC, `name` ASC LIMIT 0, 100"
	wantArgs = []interface{}{"coder", "female", "female", "coder", "hacker"}
	b.Delete("user").Where(nameEqCoder, AndSexEqFemale).Append(" /*===*/ ").Or(sexEqFemale, nameInNames).OrderBy(ageDesc, nameAsc).Limit(0, 100)
	q, err = b.Build()
	got = q.Query
	args = q.Args
	if err != nil {
		t.Errorf("delete error:%s", err)
	}

	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
	if !reflect.DeepEqual(wantArgs, args) {
		t.Errorf("\ngotArgs:\n%#v\nwantArgs:\n%#v\n", args, wantArgs)
	}
}

func TestInsert(t *testing.T) {
	var (
		got, want      string
		args, wantArgs []interface{}
		q              *Query
		err            error
	)
	// user := &Dao{
	// 	tableName: "user",
	// 	Fields:    []string{"id", "name", "age", "sex", "birthday"},
	// }
	d := &Dao{
		b:         b,
		tableName: "user",
		fields:    []string{"id", "name", "age", "sex", "birthday", "email"},
	}
	want = "INSERT INTO `user` (`id`, `name`, `age`, `sex`, `birthday`, `email`) VALUES (?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?)"
	wantArgs = []interface{}{
		1, "coder", 25, "male", "2000/09/01", "coder@coder.com",
		1, "coder", 25, "male", "2000/09/01", "coder@coder.com",
	}
	vals := []interface{}{1, "coder", 25, "male", "2000/09/01", "coder@coder.com"}
	valsGroup := [][]interface{}{vals, vals}
	b.Insert(d.TableName(), d.Fields()...).
		Values(valsGroup...)
	q, err = b.Build()
	got = q.Query
	args = q.Args

	if err != nil {
		t.Errorf("insert error:%s", err)
	}

	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
	if !reflect.DeepEqual(wantArgs, args) {
		t.Errorf("\ngotArgs:\n%#v\nwantArgs:\n%#v\n", args, wantArgs)
	}

	vals = []interface{}{1, "coder", 25, "male", "2000/09/01", "coder@coder.com"}
	valsGroup = [][]interface{}{vals}

	want = "REPLACE INTO `user` (`id`, `name`, `age`, `sex`, `birthday`, `email`) VALUES (?, ?, ?, ?, ?, ?)"
	b.Replace(d.TableName(), d.Fields()...).
		Values(valsGroup...)
	q, err = b.Build()
	got = q.Query
	args = q.Args
	if err != nil {
		t.Errorf("replace error:%s", err)
	}

	wantArgs = []interface{}{1, "coder", 25, "male", "2000/09/01", "coder@coder.com"}

	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
	if !reflect.DeepEqual(wantArgs, args) {
		t.Errorf("\ngotArgs:\n%#v\nwantArgs:\n%#v\n", args, wantArgs)
	}

	want = "INSERT INTO `user` (`id`, `name`, `age`, `sex`, `birthday`) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE `id` = ?, `name` = ?, `age` = ?, `sex` = ?, `birthday` = ?"
	wantArgs = []interface{}{1, "coder", 25, "male", "2000/09/01", 1, "coder", 25, "male", "2000/09/01"}
	vals = []interface{}{1, "coder", 25, "male", "2000/09/01"}
	valsGroup = [][]interface{}{vals, vals}
	fvals := []*FieldValue{
		NewFV("id", 1),
		NewFV("name", "coder"),
		NewFV("age", 25),
		NewFV("sex", "male"),
		NewFV("birthday", "2000/09/01"),
	}
	b.InsertOrUpdate(d.TableName(), fvals...)

	q, err = b.Build()
	got = q.Query
	args = q.Args

	if err != nil {
		t.Errorf("insert or update error:%s", err)
	}

	if want != got {
		t.Errorf("\ngot:\n{%s}\nwant:\n{%s}\n", got, want)
	}
	if !reflect.DeepEqual(wantArgs, args) {
		t.Errorf("\ngotArgs:\n%#v\nwantArgs:\n%#v\n", args, wantArgs)
	}
}

func TestReplace(t *testing.T) {
	var (
		got, want      string
		args, wantArgs []interface{}
		q              *Query
		err            error
	)

	want = "REPLACE INTO `user` (`id`, `name`, `age`, `sex`, `birthday`) VALUES (?, ?, ?, ?, ?), (?, ?, ?, ?, ?)"
	vals := []interface{}{1, "coder", 25, "male", "2000/09/01"}
	vals2 := []interface{}{2, "coder2", 25, "female", "2001/09/01"}
	valsGroup := [][]interface{}{vals}
	b.Replace(d.TableName(), d.Fields()...).
		Values(valsGroup...).
		Append(", (?, ?, ?, ?, ?)", vals2...)
	q, err = b.Build()
	got = q.Query
	args = q.Args
	if err != nil {
		t.Errorf("replace error:%s", err)
	}

	wantArgs = []interface{}{1, "coder", 25, "male", "2000/09/01", 2, "coder2", 25, "female", "2001/09/01"}

	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
	if !reflect.DeepEqual(wantArgs, args) {
		t.Errorf("\ngotArgs:\n%#v\nwantArgs:\n%#v\n", args, wantArgs)
	}
}

func TestRawBuild(t *testing.T) {
	var (
		got, want      string
		args, wantArgs []interface{}
		err            error
		q              *Query
	)
	q, err = b.Clear().Build()
	if err == nil || err != ErrEmptySQLType {
		t.Errorf("unexcept error: %#v, q=%#v", err, q)
	}
	want = "SELECT"
	q, err = b.Select().From().Build()
	got = q.Query
	args = q.Args
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}

	want = "SELECT * FROM `table` WHERE `f` = ?"
	wantArgs = []interface{}{"v"}
	b.Raw("SELECT * FROM `table` WHERE `f` = ?", "v")
	q, err = b.Build("show tables")
	got = q.Query
	args = q.Args
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
	if !reflect.DeepEqual(wantArgs, args) {
		t.Errorf("\ngotArgs:\n%#v\nwantArgs:\n%#v\n", args, wantArgs)
	}

	want = "SELECT * FROM `table` WHERE `f` = ?"
	wantArgs = []interface{}{"v"}
	b.Select("*").From("table").WhereRaw("`f` = ?", wantArgs[0])
	q, err = b.Build("show tables")
	got = q.Query
	args = q.Args
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
	if !reflect.DeepEqual(wantArgs, args) {
		t.Errorf("\ngotArgs:\n%#v\nwantArgs:\n%#v\n", args, wantArgs)
	}

	want = "show tables"
	// showTables := NewQuery("show tables", "")
	q, err = b.Raw("show tables").Build()
	got = q.Query
	args = q.Args
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}

	want = "SELECT COUNT(/*test query*/) FROM `users` WHERE `age` = ?"
	wantArgs = []interface{}{"18"}

	q, err = b.Count("/*test query*/").From("users").Where(newCondition(true, "age", "=", []interface{}{18}), nil).Build()
	got = q.Query
	args = q.Args
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}

	want = "/* just a comment */ SELECT * FROM `tablename` WHERE `name` = ?"
	wantArgs = []interface{}{"yourname"}
	b.Select("something").From("world").Clear().Select("*").From("tablename").Append(" WHERE `name` = ?", "yourname").AppendPre("/* just a comment */ ")
	q.Query = b.Query()
	q.Args = b.QueryArgs()

	q, err = b.Build("test")
	got = q.Query
	args = q.Args
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
	if !reflect.DeepEqual(wantArgs, args) {
		t.Errorf("\ngotArgs:\n%#v\nwantArgs:\n%#v\n", args, wantArgs)
	}
	lastQueries := b.LastQueries()
	lastQuery := lastQueries[len(lastQueries)-1]
	if lastQuery.Query != got {
		t.Errorf("\ngot:\n%s\nlast query:\n%s\n", got, lastQuery.Query)
	}
}

func TestSetDialector(t *testing.T) {
	var (
		got, want      string
		args, wantArgs []interface{}
		err            error
		q              *Query
	)
	t.Logf("default escape char:[%v]", b.EscapeChar())
	b.SetDialector(postgresDialector)
	t.Logf("postgres escape char:[%v]", b.EscapeChar())
	b.SetDialector(mysqlDialector)
	t.Logf("mysql escape char:[%v]", b.EscapeChar())

	want = `SELECT * FROM "user" WHERE 1 AND ("name" = ? OR "sex" = ?)`
	wantArgs = []interface{}{"coder", "female"}
	b.SetDialector(postgresDialector)

	b.Select("*").From("user").Where().And(nameEqCoder, sexEqFemale)
	q, err = b.Build()
	got = q.Query
	args = q.Args
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}

	if !reflect.DeepEqual(wantArgs, args) {
		t.Errorf("\ngotArgs:\n%#v\nwantArgs:\n%#v\n", args, wantArgs)
	}
	b.SetDialector(mysqlDialector)
}

func TestCount(t *testing.T) {
	var (
		got, want      string
		args, wantArgs []interface{}
		err            error
		q              *Query
	)
	q, err = b.Clear().Build()
	if err == nil || err != ErrEmptySQLType {
		t.Errorf("unexcept error: %#v, q=%#v", err, q)
	}
	want = "SELECT"
	q, err = b.Select().From().Build()
	got = q.Query
	args = q.Args
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}

	want = "SELECT * FROM `table` WHERE `f` = ?"
	wantArgs = []interface{}{"v"}
	b.Raw("SELECT * FROM `table` WHERE `f` = ?", "v")
	q, err = b.Build("show tables")
	got = q.Query
	args = q.Args
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
	if !reflect.DeepEqual(wantArgs, args) {
		t.Errorf("\ngotArgs:\n%#v\nwantArgs:\n%#v\n", args, wantArgs)
	}

	want = "SELECT COUNT(/*test query*/) FROM `users` WHERE `age` = ?"
	wantArgs = []interface{}{"18"}

	q, err = b.Count("/*test query*/").From("users").Where(newCondition(true, "age", "=", []interface{}{18})).Build()
	got = q.Query
	args = q.Args
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}

	want = "SELECT COUNT(1) FROM `users` WHERE `age` = ?"
	wantArgs = []interface{}{"18"}

	q, err = b.Count().From("users").Where(newCondition(true, "age", "=", []interface{}{18})).Build()
	got = q.Query
	args = q.Args
	if err != nil {
		t.Errorf("error: %s", err)
	}
	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}

	lastQueries := b.LastQueries()
	lastQuery := lastQueries[len(lastQueries)-1]
	if lastQuery.Query != got {
		t.Errorf("\ngot:\n%s\nlast query:\n%s\n", got, lastQuery.Query)
	}
}
