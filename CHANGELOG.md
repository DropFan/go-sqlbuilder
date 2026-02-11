# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.0.9] - 2026-02-11

- fix: handle nil conditions in Where() to avoid incorrect AND/OR prefixes (8a4d4d8)
- fix: handle nil FieldValues in Set() to avoid leading comma (8a4d4d8)
- fix: add nil-only conditions warning (ErrNilConditions) to ErrList in Where() (8a4d4d8)
- fix: prevent panic when Limit() called with no arguments (8a4d4d8)

## [0.0.8] - 2025-01-19

- fix: use dialect-specific placeholders for all query types (3ff275e)
- refactor: rewrite Query.String() with proper placeholder handling (f5d7789)
- feat: add query history size limits to prevent memory leaks (1d3b51b)
- feat: improve error handling with detailed error access (9b4b556)

## [0.0.7] - 2025-03-27

- performance optimization (33d5556)
- add benchmark test (7849063)

## [0.0.6] - 2025-01-28

- add missing methods in cond.go (f521f06)
- Update README.md (0aa1ca2)
- Create go.yml for GitHub Actions (ec367c6)

## [0.0.5] - 2025-01-25

- update readme, add Chinese readme (5fcb6ee)
- add SQLite dialector support (4243a8b)

## [0.0.4] - 2021-01-08

- add InsertOrUpdate(): INSERT INTO ... ON DUPLICATE UPDATE ... (3b083ec)

## [0.0.3] - 2020-09-23

- add In/NotIn/Between/NotBetween conditions & tests (96d1d94)
- go mod support (e01597c)
- update test case & comments (8edba8c)

## [0.0.2] - 2019-06-25

- replace escape char by dialector & modify some tests (bd9a911)
- add Count() method (80c5230)

## [0.0.1] - 2018-04-04

- Initial release of go-sqlbuilder
- Support for MySQL and PostgreSQL dialects
- Fluent API for SELECT, INSERT, UPDATE, DELETE queries
- add OrderBy(), Desc(), Asc() methods
- add Query.String() method for debugging
- add Condition and FieldValue types
