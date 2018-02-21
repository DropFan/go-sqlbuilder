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
	want = "SELECT `id`, `name`, `age`, `sex`, `birthday` FROM `user` WHERE [error: Invalid operator:(operator:![field:test])] AND `name` IN (?, ?) OR `sex` = ? OR `name` = ? AND [error: Invalid number of values with operator:(=[field:test_field])] ORDER BY `age` DESC, `name` ASC LIMIT 100 LIMIT 0, 100"
	wantArgs = []interface{}{"coder", "hacker", "female", "coder"}
	b.Select(d.Fields()...).
		From("user").
		Where([]Condition{errOpCond, nameInNames, sexEqFemale}...).
		Or(nameEqCoder).And(errValNumCond).
		Or().
		OrderBy(ageDesc, nameAsc).Limit(100).Limit(0, 100)
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

}

func TestUpdate(t *testing.T) {
	var (
		got, want      string
		args, wantArgs []interface{}
		q              *Query
		err            error
		fvals          = []FieldValue{
			FieldValue{Name: "tag", Value: "test"},
			FieldValue{Name: "desc", Value: "just 4 test"},
		}
		fv = FieldValue{
			Name:  "f",
			Value: "v",
		}
		kv = FieldValue{
			Name:  "k",
			Value: "v",
		}
	)
	AndSexEqFemale := sexEqFemale
	AndSexEqFemale.AndOr = true

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
}

func TestDelete(t *testing.T) {

	var (
		got, want      string
		args, wantArgs []interface{}
		q              *Query
		err            error
	)
	AndSexEqFemale := sexEqFemale
	AndSexEqFemale.AndOr = true

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
	want = "INSERT INTO `user` (`id`, `name`, `age`, `sex`, `birthday`) VALUES (?, ?, ?, ?, ?)"
	vals := []interface{}{1, "coder", 25, "male", "2000/09/01"}
	valsGroup := [][]interface{}{vals}
	b.Insert(d.TableName(), d.Fields()...).
		Values(valsGroup...)
	q, err = b.Build()
	got = q.Query
	args = q.Args

	if err != nil {
		t.Errorf("insert error:%s", err)
	}

	wantArgs = []interface{}{1, "coder", 25, "male", "2000/09/01"}

	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
	}
	if !reflect.DeepEqual(wantArgs, args) {
		t.Errorf("\ngotArgs:\n%#v\nwantArgs:\n%#v\n", args, wantArgs)
	}

	vals = []interface{}{1, "coder", 25, "male", "2000/09/01"}
	valsGroup = [][]interface{}{vals}

	want = "REPLACE INTO `user` (`id`, `name`, `age`, `sex`, `birthday`) VALUES (?, ?, ?, ?, ?)"
	b.Replace(d.TableName(), d.Fields()...).
		Values(valsGroup...)
	q, err = b.Build()
	got = q.Query
	args = q.Args
	if err != nil {
		t.Errorf("replace error:%s", err)
	}

	wantArgs = []interface{}{1, "coder", 25, "male", "2000/09/01"}

	if want != got {
		t.Errorf("\ngot:\n%s\nwant:\n%s\n", got, want)
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
	valsGroup := [][]interface{}{vals}
	b.Replace(d.TableName(), d.Fields()...).
		Values(valsGroup...).
		Append(", (?, ?, ?, ?, ?)", vals...)
	q, err = b.Build()
	got = q.Query
	args = q.Args
	if err != nil {
		t.Errorf("replace error:%s", err)
	}

	wantArgs = []interface{}{1, "coder", 25, "male", "2000/09/01", 1, "coder", 25, "male", "2000/09/01"}

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

	want = "/* just a comment */ SELECT * FROM `tablename` WHERE `name` = ?"
	wantArgs = []interface{}{"yourname"}
	b.Select("*").From("tablename").Append(" WHERE `name` = ?", "yourname").AppendPre("/* just a comment */ ")
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
