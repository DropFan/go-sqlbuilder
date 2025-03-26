package builder

import (
	"testing"
)

// Common setup for benchmark tests
var (
	benchBuilder = New()
	benchFields  = []string{"id", "name", "age", "sex", "birthday"}
	benchTable   = "user"
)

// BenchmarkSelect tests the performance of basic Select query building
func BenchmarkSelect(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchBuilder.Select(benchFields...).From(benchTable).Build()
	}
}

// BenchmarkSelectWithWhere tests the performance of Select query building with Where conditions
func BenchmarkSelectWithWhere(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchBuilder.Select(benchFields...).From(benchTable).
			Where(Eq("age", 25)).And(Eq("sex", "male")).
			Build()
	}
}

// BenchmarkSelectComplex tests the performance of complex Select query building
func BenchmarkSelectComplex(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchBuilder.Select(benchFields...).From(benchTable).
			Where(Gt("age", 18)).
			And(In("city_id", 1, 2, 3)).
			Or(Between("salary", 5000, 10000)).
			OrderBy(Desc("age"), Asc("name")).
			Limit(0, 100).
			Build()
	}
}

// BenchmarkInsert tests the performance of basic Insert query building
func BenchmarkInsert(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchBuilder.Insert(benchTable, benchFields...).
			Values([]interface{}{1, "John", 25, "male", "1998-01-01"}).
			Build()
	}
}

// BenchmarkInsertMultipleRows tests the performance of Insert query building with multiple rows
func BenchmarkInsertMultipleRows(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	values1 := []interface{}{1, "John", 25, "male", "1998-01-01"}
	values2 := []interface{}{2, "Jane", 23, "female", "2000-01-01"}
	values3 := []interface{}{3, "Bob", 30, "male", "1993-01-01"}

	for i := 0; i < b.N; i++ {
		benchBuilder.Insert(benchTable, benchFields...).
			Values(values1, values2, values3).
			Build()
	}
}

// BenchmarkUpdate tests the performance of basic Update query building
func BenchmarkUpdate(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchBuilder.Update(benchTable,
			NewFieldValue("name", "John Doe"),
			NewFieldValue("age", 26),
		).Where(Eq("id", 1)).Build()
	}
}

// BenchmarkUpdateComplex tests the performance of complex Update query building
func BenchmarkUpdateComplex(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchBuilder.Update(benchTable,
			NewFieldValue("name", "John Doe"),
			NewFieldValue("age", 26),
			NewFieldValue("updated_at", "2023-01-01"),
		).Where(Eq("id", 1)).
			And(Gt("age", 20)).
			Or(In("status", "active", "pending")).
			Build()
	}
}

// BenchmarkDelete tests the performance of basic Delete query building
func BenchmarkDelete(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchBuilder.Delete(benchTable).Where(Eq("id", 1)).Build()
	}
}

// BenchmarkDeleteComplex tests the performance of complex Delete query building
func BenchmarkDeleteComplex(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchBuilder.Delete(benchTable).
			Where(Lt("created_at", "2020-01-01")).
			And(Eq("status", "inactive")).
			Or(In("category", "temp", "test")).
			Build()
	}
}

// BenchmarkRawSQL tests the performance of raw SQL building
func BenchmarkRawSQL(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchBuilder.Raw("SELECT * FROM user WHERE id = ?", 1).Build()
	}
}

// BenchmarkDialectorMysql tests the performance of MySQL dialect
func BenchmarkDialectorMysql(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	builder := New()
	builder.SetDialector(mysqlDialector)
	for i := 0; i < b.N; i++ {
		builder.Select(benchFields...).From(benchTable).
			Where(Eq("id", 1)).
			Build()
	}
}

// BenchmarkDialectorPostgresql tests the performance of PostgreSQL dialect
func BenchmarkDialectorPostgresql(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	builder := New()
	builder.SetDialector(postgresDialector)
	for i := 0; i < b.N; i++ {
		builder.Select(benchFields...).From(benchTable).
			Where(Eq("id", 1)).
			Build()
	}
}

// BenchmarkDialectorSQLite tests the performance of SQLite dialect
func BenchmarkDialectorSQLite(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	builder := New()
	builder.SetDialector(sqliteDialector)
	for i := 0; i < b.N; i++ {
		builder.Select(benchFields...).From(benchTable).
			Where(Eq("id", 1)).
			Build()
	}
}

// BenchmarkConditionBuilding tests the performance of condition building
func BenchmarkConditionBuilding(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Eq("id", 1)
		Gt("age", 18)
		Lt("age", 60)
		Gte("salary", 5000)
		Lte("salary", 10000)
		In("status", "active", "pending", "reviewing")
		NotIn("category", "deleted", "archived")
		Between("created_at", "2020-01-01", "2023-01-01")
		NotBetween("updated_at", "2022-01-01", "2022-12-31")
		Like("name", "%John%")
		NotLike("description", "%temp%")
	}
}

// BenchmarkComplexConditionCombination tests the performance of complex condition combinations
func BenchmarkComplexConditionCombination(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		condGroup := NewConditionGroup(
			Eq("status", "active"),
			Gt("age", 18),
		)
		condGroup2 := NewConditionGroup(
			In("category", "premium", "standard"),
			Between("created_at", "2020-01-01", "2023-01-01"),
		)

		benchBuilder.Select(benchFields...).From(benchTable).
			Where(condGroup...).
			Or(condGroup2...).
			Build()
	}
}
